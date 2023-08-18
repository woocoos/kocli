package appres

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
	_ "github.com/woocoos/knockout/ent/runtime"
)

var (
	cnf = &cli.PathFlag{
		Name:    "config",
		Usage:   "the knockout config",
		Value:   "knockout.yaml",
		Aliases: []string{"f"},
	}
	appcode = &cli.StringFlag{
		Name:    "app",
		Usage:   "application code",
		Aliases: []string{"a"},
	}
)

var Cmd = &cli.Command{
	Name:  "res",
	Usage: "tools for knockout resource generator",
	Subcommands: []*cli.Command{
		{
			Name:  "gql-action",
			Usage: "gen actions from graphql",
			Flags: []cli.Flag{
				cnf,
				appcode,
				&cli.PathFlag{
					Name:    "gql",
					Usage:   "the gqlgen config path",
					Value:   "codegen/gqlgen/gqlgen.yaml",
					Aliases: []string{"g"},
				},
			},
			Action: func(ctx *cli.Context) error {
				return GenGqlActions(Config{
					KnockoutConfig: ctx.String("config"),
					GQLConfig:      ctx.String("gql"),
					AppCode:        ctx.String("app"),
				})
			},
		},
		{
			Name:  "ent",
			Usage: "gen resource from ent schema. Root schema will generate to app_res",
			Flags: []cli.Flag{
				cnf,
				appcode,
				&cli.PathFlag{
					Name:    "schema",
					Usage:   "the ent schema path",
					Value:   "codegen/entgen/schema",
					Aliases: []string{"e"},
				},
			},
			Action: func(ctx *cli.Context) error {
				return GenEntSchemaRes(Config{
					KnockoutConfig: ctx.String("config"),
					EntConfig:      ctx.String("schema"),
					AppCode:        ctx.String("app"),
				})
			},
		},
		{
			Name:  "menu",
			Usage: "gen app menu from web project.",
			Flags: []cli.Flag{
				cnf,
				appcode,
				&cli.PathFlag{
					Name:    "data",
					Usage:   "the menu data path",
					Value:   "web/src/components/Layout/menu.json",
					Aliases: []string{"d"},
				},
			},
			Action: func(ctx *cli.Context) error {
				return GenAppMenu(Config{
					KnockoutConfig: ctx.String("config"),
					EntConfig:      ctx.String("schema"),
					MenuConfig:     ctx.String("data"),
					AppCode:        ctx.String("app"),
				})
			},
		},
	},
}
