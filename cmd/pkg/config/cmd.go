package config

import (
	"github.com/4396/pkg/config"
	"github.com/urfave/cli"
)

var (
	token    string
	registry string

	Command = cli.Command{
		Name:  "config",
		Usage: "set registry or token",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "token",
				Usage:       "",
				Destination: &token,
			},
			cli.StringFlag{
				Name:        "registry",
				Usage:       "",
				Destination: &registry,
			},
		},
		Action: func(c *cli.Context) (err error) {
			_, err = config.Registry(registry)
			if err != nil {
				return
			}

			_, err = config.Token(token)
			return
		},
	}
)
