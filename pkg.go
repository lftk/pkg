package main

import (
	"log"
	"os"

	"github.com/4396/pkg/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "pkg"
	app.Usage = "used to download packages and dependencies"
	app.Version = "v0.1.0"
	app.Commands = []cli.Command{
		cmd.Get,
		cmd.Mirror,
		cmd.Config,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
