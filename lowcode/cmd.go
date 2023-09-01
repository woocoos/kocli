package lowcode

import (
	"github.com/urfave/cli/v2"
	"github.com/woocoos/kocli/lowcode/graph"
	"path/filepath"
)

var PageCmd = &cli.Command{
	Name:  "page",
	Usage: "a tool for generate knockout application",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "target",
			Aliases: []string{"t"},
			Usage:   "the target directory of the generated code",
			Value:   ".",
		},
	},
	Action: func(c *cli.Context) (err error) {
		dir := c.String("target")
		fd, err := filepath.Abs(dir)
		if err != nil {
			return err
		}
		cfg := &graph.Config{
			Target: fd,
		}
		var opts []Option
		return Generate(cfg, opts...)
	},
}
