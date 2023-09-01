package project

import (
	"bytes"
	"github.com/tsingsun/woocoo/cmd/woco/gen"
	"github.com/tsingsun/woocoo/cmd/woco/project"
	"path/filepath"
)

var _ gen.Extension = (*Extension)(nil)

var (
	mainT = gen.ParseT("template/main.tmpl", templateDir, ExtensionFuncs)
)

type Option func(*Extension)

func WithTargetDir(dir string) Option {
	return func(e *Extension) {
		e.TargetDir = dir
	}
}

// WithSkipRunGen option for debug and testing, when run entgen and gqlgen in testing will failure
func WithSkipRunGen(skip bool) Option {
	return func(e *Extension) {
		e.SkipRunGen = skip
	}
}

func WithFrontend() Option {
	return func(e *Extension) {
		t := gen.NewTemplate("frontend")
		t.Delims("[[", "]]")
		e.templates = append(e.templates,
			gen.MustParse(t.Funcs(ExtensionFuncs).ParseFS(templateDir,
				"template/web/*.tmpl", "template/web/**/*.tmpl", "template/web/**/**/*.tmpl")),
		)
		staticFiles := map[string]string{
			"template/web/public/favicon.ico": "web/public/favicon.ico",
		}
		for k, v := range staticFiles {
			e.staticFiles[k] = v
		}
	}
}

type Extension struct {
	TargetDir string
	// SkipRunGen option for debug and testing, when run entgen and gqlgen in testing will failure
	SkipRunGen  bool
	hooks       []gen.Hook
	templates   []*gen.Template
	staticFiles map[string]string
}

func New(opt ...Option) *Extension {
	ex := &Extension{
		staticFiles: make(map[string]string),
	}
	ex.initTemplates()
	for _, o := range opt {
		o(ex)
	}
	return ex
}

func (e *Extension) writeStatic() error {
	var assets gen.Assets
	for key, tar := range e.staticFiles {
		fp := filepath.Join(e.TargetDir, tar)
		assets.AddDir(filepath.Dir(fp))
		bs, err := templateDir.ReadFile(key)
		if err != nil {
			return err
		}
		assets.Add(fp, bs)
	}
	return assets.Write()
}

func (e *Extension) GeneratedHooks() []gen.GeneratedHook {
	return []gen.GeneratedHook{
		func(ex gen.Extension) error {
			if err := e.writeStatic(); err != nil {
				return err
			}

			graph := ex.(*project.Graph)
			if e.SkipRunGen {
				return nil
			}
			if graph.SkipModTidy {
				_ = runGoModTidy(e.TargetDir)
			}
			err := gen.RunCmd(e.TargetDir, "go", "run", "codegen/entgen/entc.go")
			if err != nil {
				return err
			}
			_ = runGoModTidy(e.TargetDir)
			// ignore gqlgen error
			_ = gen.RunCmd(e.TargetDir, "go", "run", "codegen/gqlgen/gqlgen.go")

			b := &bytes.Buffer{}
			err = mainT.ExecuteTemplate(b, "main", graph)
			if err != nil {
				return err
			}
			if err = gen.FormatGoFile(filepath.Join(e.TargetDir, "cmd/main.go"), b.Bytes()); err != nil {
				return err
			}
			_ = runGoModTidy(e.TargetDir)

			return nil
		},
	}
}

func (e *Extension) Name() string {
	return "knockout"
}

func (e *Extension) Templates() []*gen.Template {
	return e.templates
}

// Hooks knockout not need write file itself
func (e *Extension) Hooks() []gen.Hook {
	return e.hooks
}

func (e *Extension) initTemplates() {
	initTemplates()
	e.templates = append(e.templates,
		gen.ParseT("template/cmd.tmpl", templateDir, ExtensionFuncs),
		gen.ParseT("template/codegen.tmpl", templateDir, ExtensionFuncs),
		gen.ParseT("template/makefile.tmpl", templateDir, ExtensionFuncs),
		gen.ParseT("template/test.tmpl", templateDir, ExtensionFuncs),
		//gen.ParseT("template/main.tmpl", templateDir, ExtensionFuncs),
		gen.ParseT("template/gomod.tmpl", templateDir, ExtensionFuncs),
	)
}

func runGoModTidy(dir string) error {
	return gen.RunCmd(dir, "go", "mod", "tidy")
}
