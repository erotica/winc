package main

import (
	"code.cloudfoundry.org/winc/hcs"
	"code.cloudfoundry.org/winc/sandbox"
	"github.com/urfave/cli"
)

var statsCommand = cli.Command{
	Name:      "stats",
	Usage:     "show stats for container volume",
	ArgsUsage: `<container-id>`,
	Action: func(context *cli.Context) error {
		if err := checkArgs(context, 1, exactArgs); err != nil {
			return err
		}

		containerId := context.Args().First()
		storePath := context.GlobalString("store")

		sm := sandbox.NewManager(&hcs.Client{}, nil, storePath, containerId)
		return sm.Stats()
	},
}
