package project

import (
	"github.com/tsingsun/woocoo/cmd/woco/project"
	"github.com/urfave/cli/v2"
	"path/filepath"
)

var InitCmd = &cli.Command{
	Name:  "init",
	Usage: "a tool for generate knockout application",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "package",
			Aliases: []string{"p"},
			Usage:   "the package name of the generated code",
			Value:   "github.com/woocoos/helloworld",
		},
		&cli.StringFlag{
			Name:    "target",
			Aliases: []string{"t"},
			Usage:   "the target directory of the generated code",
			Value:   ".",
		},
		&cli.BoolFlag{
			Name:    "frontend",
			Aliases: []string{"f"},
			Usage:   "generate frontend code",
		},
	},
	Action: func(c *cli.Context) (err error) {
		dir := c.String("target")
		// get full path by "."
		fd, err := filepath.Abs(dir)
		if err != nil {
			return err
		}
		cfg := &project.Config{
			Package:     c.String("package"),
			Target:      fd,
			Modules:     []string{"otel", "web"},
			SkipModTidy: true,
		}
		exts := []Option{
			WithTargetDir(cfg.Target),
		}
		if c.Bool("frontend") {
			exts = append(exts, WithFrontend())
		}

		var opts []project.Option
		opts = append(opts, project.Extensions(New(exts...)))
		return project.Generate(cfg, opts...)
	},
}
