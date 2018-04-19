package cmd

import (
	"fmt"

	"github.com/4396/pkg/mirror"
	"github.com/urfave/cli"
)

var (
	cmdSet = cli.Command{
		Name:  "set",
		Usage: "set mirror",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "pkg",
			},
			cli.StringFlag{
				Name: "repo",
			},
			cli.StringFlag{
				Name: "base",
			},
		},
		Action: func(c *cli.Context) error {
			pkg := c.String("pkg")
			repo := c.String("repo")
			base := c.String("base")
			if pkg == "" || repo == "" || base == "" {
				return nil
			}
			return mirror.Set(pkg, repo, base)
		},
	}

	cmdGet = cli.Command{
		Name:  "get",
		Usage: "get mirror",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "pkg",
			},
		},
		Action: func(c *cli.Context) error {
			pkg := c.String("pkg")
			if pkg == "" {
				return nil
			}

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
				Name: "pkg",
			},
		},
		Action: func(c *cli.Context) error {
			pkg := c.String("pkg")
			if pkg == "" {
				return nil
			}
			return mirror.Delete(pkg)
		},
	}

	Mirror = cli.Command{
		Name:  "mirror",
		Usage: "forwarding package",
		Subcommands: []cli.Command{
			cmdSet, cmdGet, cmdDel,
		},
	}
)
