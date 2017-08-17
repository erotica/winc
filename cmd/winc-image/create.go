package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"code.cloudfoundry.org/winc/hcs"
	"code.cloudfoundry.org/winc/image"
	"code.cloudfoundry.org/winc/layer"
	"code.cloudfoundry.org/winc/volume"

	"github.com/urfave/cli"
)

var createCommand = cli.Command{
	Name:      "create",
	Usage:     "create a container volume",
	ArgsUsage: `<rootfs> <container-id>`,
	Flags: []cli.Flag{
		cli.Uint64Flag{
			Name:  "disk-limit-size-bytes",
			Usage: "Disk limit in bytes",
		},
		cli.BoolFlag{
			Name:  "exclude-image-from-quota",
			Usage: "Ignored",
		},
	},
	Action: func(context *cli.Context) error {
		if err := checkArgs(context, 2, exactArgs); err != nil {
			return err
		}

		rootfsPath := context.Args().First()
		containerId := context.Args().Tail()[0]
		storePath := context.GlobalString("store")
		diskLimit := context.Uint64("disk-limit-size-bytes")

		destToWindowsPath := func(input string) string {
			vol := filepath.VolumeName(input)
			if vol == "" {
				input = filepath.Join("C:", input)
			}
			return filepath.Clean(input)
		}

		rootfsPath = destToWindowsPath(rootfsPath)
		lm := layer.NewManager(&hcs.Client{}, storePath)
		im := image.NewManager(lm, &volume.Limiter{}, &volume.Statser{}, containerId)
		imageSpec, err := im.Create(rootfsPath, diskLimit)
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
