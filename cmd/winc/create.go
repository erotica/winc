package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"unicode/utf8"

	"code.cloudfoundry.org/winc/container"
	"code.cloudfoundry.org/winc/hcsclient"
	"code.cloudfoundry.org/winc/sandbox"

	"github.com/Sirupsen/logrus"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/opencontainers/runtime-tools/validate"
	"github.com/urfave/cli"
)

var createCommand = cli.Command{
	Name:  "create",
	Usage: "create a container",
	ArgsUsage: `<container-id>

Where "<container-id>" is your name for the instance of the container that you
are starting. The name you provide for the container instance must be unique on
your host.`,
	Description: `The create command creates an instance of a container for a bundle. The bundle
	is a directory with a specification file named "` + specConfig + `" and a root
	filesystem.

	The specification file includes an args parameter. The args parameter is used
	to specify command(s) that get run when the container is started. To change the
	command(s) that get executed on start, edit the args parameter of the spec`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "bundle, b",
			Value: "",
			Usage: `path to the root of the bundle directory, defaults to the current directory`,
		},
		cli.StringFlag{
			Name:  "pid-file",
			Value: "",
			Usage: "specify the file to write the process id to",
		},
		cli.BoolFlag{
			Name:  "no-new-keyring",
			Usage: "ignored",
		},
	},
	Action: func(context *cli.Context) error {
		if err := checkArgs(context, 1, exactArgs); err != nil {
			return err
		}

		containerId := context.Args().First()
		bundlePath := context.String("bundle")
		pidFile := context.String("pid-file")

		logrus.WithFields(logrus.Fields{
			"bundle":      bundlePath,
			"containerId": containerId,
			"pidFile":     pidFile,
		}).Debug("creating container")

		if bundlePath == "" {
			var err error
			bundlePath, err = os.Getwd()
			if err != nil {
				return err
			}
		}

		if _, err := os.Stat(bundlePath); err != nil {
			return &MissingBundleError{BundlePath: bundlePath}
		}

		configPath := filepath.Join(bundlePath, specConfig)
		content, err := ioutil.ReadFile(configPath)
		if err != nil {
			return &MissingBundleConfigError{BundlePath: bundlePath}
		}
		if !utf8.Valid(content) {
			return fmt.Errorf("%q is not encoded in UTF-8", configPath)
		}
		var spec specs.Spec
		if err = json.Unmarshal(content, &spec); err != nil {
			return &BundleConfigInvalidJSONError{}
		}

		validator := validate.NewValidator(&spec, bundlePath, true)

		m := validator.CheckMandatoryFields()
		if len(m) != 0 {
			return &BundleConfigValidationError{m}
		}

		client := hcsclient.HCSClient{}
		sm := sandbox.NewManager(&client, bundlePath)
		cm := container.NewManager(&client, sm, containerId)

		if err := cm.Create(spec.Root.Path); err != nil {
			return err
		}

		if pidFile != "" {
			state, err := cm.State()
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(pidFile, []byte(strconv.FormatInt(int64(state.Pid), 10)), 0666); err != nil {
				return err
			}
		}

		return nil
	},
}
