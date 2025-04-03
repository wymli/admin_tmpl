package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "xserver",
		Usage: "xserver tools",
		Commands: []*cli.Command{
			{
				Name:    "generate",
				Aliases: []string{"gen"},
				Usage:   "generate tools",
				Commands: []*cli.Command{{
					Name:    "view",
					Aliases: nil,
					Usage:   "generate view routes and handlers",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "src",
							Value:    "",
							Usage:    "src dir for view models",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "dest",
							Value:    "",
							Usage:    "dest dir for view handler",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "mod",
							Value:    "",
							Usage:    "go module path",
							Required: true,
						},
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						g := &ViewGenerator{
							SrcDir:  cmd.String("src"),
							DestDir: cmd.String("dest"),
							ModPath: cmd.String("mod"),
						}

						if err := g.Generate(); err != nil {
							panic(err)
						}
						return nil
					},
				}},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}
