package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

var specCommand = cli.Command{
	Name:        "spec",
	Usage:       "create a new specification file",
	Description: `The spec command creates the new specification file named "config.json" for the bundle.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "bundle, b",
			Value: "",
			Usage: "path to the root of the bundle directory",
		},
	},
	Action: func(context *cli.Context) error {
		bundle := context.String("bundle")
		if bundle != "" {
			if err := os.Chdir(bundle); err != nil {
				return err
			}
		}

		rootfsPath, err := exec.Command("powershell.exe", "(docker inspect cloudfoundry/cfwindowsfs | ConvertFrom-Json).GraphDriver.Data.Dir").CombinedOutput()
		if err != nil {
			return err
		}

		spec := specs.Spec{
			Version: specs.Version,
			Platform: specs.Platform{
				OS:   runtime.GOOS,
				Arch: runtime.GOARCH,
			},
			Process: &specs.Process{
				Args: []string{"powershell"},
				Cwd:  "/",
			},
			Root: specs.Root{
				Path: strings.TrimSpace(string(rootfsPath)),
			},
		}

		config, err := json.Marshal(spec)
		if err != nil {
			return err
		}

		return ioutil.WriteFile("config.json", config, 0644)
	},
}
