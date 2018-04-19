package cmd

import (
	"github.com/4396/pkg/conf"
	"github.com/urfave/cli"
)

var (
	Config = cli.Command{
		Name:  "config",
		Usage: "set registry or token",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "token",
				Usage: "config token",
			},
			cli.StringFlag{
				Name:  "registry",
				Usage: "config registry",
			},
		},
		Action: func(c *cli.Context) (err error) {
			for name, f := range map[string]func(string) (string, error){
				"token":    conf.Token,
				"registry": conf.Registry,
			} {
				val := c.String(name)
				if val == "" {
					continue
				}

				_, err = f(val)
				if err != nil {
					return
				}
			}
			return
		},
	}
)
