package main

import (
	"os"

	"github.com/urfave/cli"
)

const (
	usage = `FIXME`
)

func main() {
	app := cli.NewApp()
	app.Name = "winc-image.exe"
	app.Usage = usage

	app.Commands = []cli.Command{
		createCommand,
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
