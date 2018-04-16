package main

import (
	"log"
	"os"

	"github.com/4396/pkg/cmd/pkg/config"
	"github.com/4396/pkg/cmd/pkg/get"
	"github.com/4396/pkg/cmd/pkg/mirror"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "pkg"
	app.Usage = "used to download packages and dependencies"
	app.Version = "v0.1.0"
	app.Commands = []cli.Command{
		get.Command,
		mirror.Command,
		config.Command,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
