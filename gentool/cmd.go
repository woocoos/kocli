package gentool

import (
	"github.com/urfave/cli/v2"
)

var ToolCmd = &cli.Command{
	Name:  "gen",
	Usage: "generate tools",
	Subcommands: []*cli.Command{
		enumCMD,
	},
}

var enumCMD = &cli.Command{
	Name:  "enum",
	Usage: "generate enum for schema",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "baseType",
			Usage:    "the base type for enum, like int or string",
			Aliases:  []string{"b"},
			Required: true,
			Value:    "int",
		},
		&cli.StringFlag{
			Name:     "name",
			Usage:    "the enum name",
			Aliases:  []string{"n"},
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "values",
			Usage:    "the enum values",
			Aliases:  []string{"v"},
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		return GenerateEnum(EnumInput{
			BaseType: c.String("baseType"),
			EnumName: c.String("name"),
			Values:   c.StringSlice("values"),
		})
	},
}
