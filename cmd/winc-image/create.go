package main

import (
	"encoding/json"
	"fmt"

	"code.cloudfoundry.org/winc/hcsclient"
	"code.cloudfoundry.org/winc/mounter"
	"code.cloudfoundry.org/winc/sandbox"

	"github.com/urfave/cli"
)

var createCommand = cli.Command{
	Name:      "create",
	Usage:     "FIXME",
	ArgsUsage: `FIXME`,
	Action: func(context *cli.Context) error {
		rootfsPath := context.Args().First()
		containerId := context.Args().Tail()[0]

		sm := sandbox.NewManager(&hcsclient.HCSClient{}, &mounter.Mounter{}, containerId)
		imageSpec, err := sm.Create(rootfsPath)
		if err != nil {
			return err
		}

		output, err := json.Marshal(&imageSpec)
		if err != nil {
			return err
		}

		fmt.Println(string(output))

		return nil
	},
}
