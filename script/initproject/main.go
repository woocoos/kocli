//go:build ignore

package main

import (
	"flag"
	"github.com/tsingsun/woocoo/cmd/woco/project"
	exprj "github.com/woocoos/kocli/project"
	"log"
)

var (
	pkg    = flag.String("package", "github.com/woocoos/helloworld", "package name")
	target = flag.String("target", ".", "target dir")
)

func main() {
	flag.Parse()
	var opts []project.Option

	cfg := &project.Config{
		Package: *pkg,
		Target:  *target,
		Modules: []string{"otel", "web"},
	}

	opts = append(opts, project.Extensions(exprj.New(exprj.WithTargetDir(cfg.Target))))
	err := project.Generate(cfg, opts...)
	if err != nil {
		log.Panic(err)
	}
}
