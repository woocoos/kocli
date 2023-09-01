package lowcode

import (
	"github.com/tsingsun/woocoo/cmd/woco/gen"
	"github.com/woocoos/kocli/lowcode/graph"
	"github.com/woocoos/kocli/lowcode/schema"
	"os"
	"path/filepath"
)

type Option func(*graph.Config) error

func Extensions(extensions ...gen.Extension) Option {
	return func(config *graph.Config) error {
		for _, ex := range extensions {
			config.Hooks = append(config.Hooks, ex.Hooks()...)
			config.Templates = append(config.Templates, ex.Templates()...)
			config.GeneratedHooks = append(config.GeneratedHooks, ex.GeneratedHooks()...)
		}
		return nil
	}
}

func Generate(cfg *graph.Config, opts ...Option) error {
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return err
		}
	}
	bs, err := os.ReadFile(cfg.Schema)
	if err != nil {
		return err
	}
	sch, err := schema.ParseProjectSchema(bs)
	if err != nil {
		return err
	}
	cfg.GeneratedHooks = append(cfg.GeneratedHooks, func(extension gen.Extension) error {
		// run npm install
		// cfg.Schema.NpmInstall()
		//return gen.RunCmd(cfg.Target, "go", "mod", "tidy")
		return nil
	})

	if cfg.Target == "" {
		abs, err := filepath.Abs(cfg.Schema)
		if err != nil {
			return err
		}
		cfg.Target = filepath.Dir(abs)
	}
	lg, err := graph.NewGraph(cfg, sch)
	if err != nil {
		return err
	}
	return lg.Gen()
}
