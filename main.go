package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version",
		Usage: "show version",
	}
	app := &cli.App{
		Name:  "verifyenv",
		Usage: "file and Environment verify tool.\nUsualy usecase are file exist check and environment value exist check.",
		Authors: []*cli.Author{{
			Name:  "Kenichiro Gozu",
			Email: "gozuk16@gmail.com",
		}},
		Version: "0.1.0",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "expand",
				Aliases: []string{"e"},
				Usage:   "Be expand when checking for duplicates, in clude other directories.",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Be verbose when validating, showing them as they are validated.",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
		Action: func(cCtx *cli.Context) error {
			optE := cCtx.Bool("expand")
			optV := cCtx.Bool("verbose")
			config := cCtx.String("config")
			fmt.Println("expand:", optE)
			fmt.Println("verbose:", optV)
			fmt.Println("config:", config)

			if config != "" {
				validateFiles(config, optE, optV)
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "add a file list from given directory.",
				Action: func(cCtx *cli.Context) error {
					if cCtx.NArg() > 0 {
						//validateFiles(config, optE, optV)
						searchFile(cCtx.Args().First(), optA)
					}
					return nil
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "all",
						Aliases: []string{"a"},
						Usage:   "Include directory entries whose names begin with a dot (.).",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
