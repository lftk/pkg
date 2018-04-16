package mirror

import (
	"fmt"

	"github.com/4396/pkg/mirror"
	"github.com/urfave/cli"
)

var (
	pkg  string
	repo string
	base string

	cmdSet = cli.Command{
		Name:  "set",
		Usage: "set mirror",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "pkg",
				Destination: &pkg,
			},
			cli.StringFlag{
				Name:        "repo",
				Destination: &repo,
			},
			cli.StringFlag{
				Name:        "base",
				Destination: &base,
			},
		},
		Action: func(c *cli.Context) error {
			return mirror.Set(pkg, repo, base)
		},
	}

	cmdGet = cli.Command{
		Name:  "get",
		Usage: "get mirror",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "pkg",
				Destination: &pkg,
			},
		},
		Action: func(c *cli.Context) error {
			repo, base, ok := mirror.Get(pkg)
			if ok {
				fmt.Println(repo, base)
			}
			return nil
		},
	}

	cmdDel = cli.Command{
		Name:  "del",
		Usage: "delete mirror",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "pkg",
				Destination: &pkg,
			},
		},
		Action: func(c *cli.Context) error {
			mirror.Delete(pkg)
			return nil
		},
	}

	Command = cli.Command{
		Name:  "mirror",
		Usage: "forwarding package",
		Subcommands: []cli.Command{
			cmdSet, cmdGet, cmdDel,
		},
	}
)
