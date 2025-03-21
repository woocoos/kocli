package main

import (
	"github.com/urfave/cli/v2"
	"github.com/woocoos/kocli/appres"
	"github.com/woocoos/kocli/gentool"
	"github.com/woocoos/kocli/project"
	"log"
	"os"
)

const Version = "0.0.1"

var commands = []*cli.Command{
	project.InitCmd,
	appres.Cmd,
	gentool.ToolCmd,
}

func main() {
	app := cli.NewApp()
	app.Name = "kocli"
	app.Usage = "a cli command for woocoos application"
	app.Version = Version
	app.Commands = commands
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err.Error())
	}
}
