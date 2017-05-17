package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/winc/container"

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
	},
	Action: func(context *cli.Context) error {
		if err := checkArgs(context, 1, exactArgs); err != nil {
			return err
		}

		bundlePath := context.String("bundle")
		containerId := context.Args().First()

		if bundlePath == "" {
			var err error
			bundlePath, err = os.Getwd()
			if err != nil {
				return err
			}
		}

		v, err := validate.NewValidatorFromPath(bundlePath, true)
		if err != nil {
			return err
		}

		m := v.CheckMandatoryFields()
		if len(m) != 0 {
			return &WincBundleConfigValidationError{m}
		}

		configBytes, err := ioutil.ReadFile(filepath.Join(bundlePath, specConfig))
		if err != nil {
			return err
		}

		type spec struct {
			Root struct {
				Path string
			}
		}

		var s spec
		if err := json.Unmarshal(configBytes, &s); err != nil {
			return err
		}

		return container.Create(s.Root.Path, bundlePath, containerId)
	},
}
