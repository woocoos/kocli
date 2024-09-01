package lowcode

import (
	"github.com/urfave/cli/v2"
	"github.com/woocoos/kocli/lowcode/graph"
	"path/filepath"
)

var PageCmd = &cli.Command{
	Name:  "page",
	Usage: "a tool for generate react web page",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "target",
			Aliases:  []string{"t"},
			Usage:    "the target directory of the generated code",
			Value:    ".",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "schema",
			Aliases:  []string{"s"},
			Usage:    "the schema file of ali lowcode protocol",
			Required: true,
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
