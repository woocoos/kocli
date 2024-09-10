package project

import (
	"bytes"
	"github.com/tsingsun/woocoo/cmd/woco/gen"
	"github.com/tsingsun/woocoo/cmd/woco/project"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var _ gen.Extension = (*Extension)(nil)

var (
	mainT = gen.ParseT("entry", templateDir, ExtensionFuncs,
		"template/main.tmpl", "template/graphql.tmpl")
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

// WithFrontend option for generate frontend code
func WithFrontend() Option {
	return func(e *Extension) {
		e.genWeb = true
		t := gen.NewTemplate("frontend")
		t.Delims("[[", "]]")
		e.templates = append(e.templates,
			gen.MustParse(t.Funcs(ExtensionFuncs).ParseFS(templateDir,
				"template/web/*.tmpl", "template/web/**/*.tmpl", "template/web/**/**/*.tmpl")),
		)
	}
}

// Extension is the knockout extension for the project.
type Extension struct {
	TargetDir  string
	SkipRunGen bool
	genWeb     bool
	hooks      []gen.Hook
	templates  []*gen.Template
}

// GeneratedHooks return the generated hooks.
// The generated hooks are executed after the code generation.
func (e *Extension) GeneratedHooks() []gen.GeneratedHook {
	return []gen.GeneratedHook{
		e.codegen,
		e.webStaticAssets,
	}
}

// codegen run entgen and gqlgen
func (e *Extension) codegen(ex gen.Extension) error {
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
	err = runGoModTidy(e.TargetDir)
	if err != nil {
		log.Println(err)
	}
	// ignore gqlgen error
	err = gen.RunCmd(e.TargetDir, "go", "run", "codegen/gqlgen/gqlgen.go")
	if err != nil {
		log.Println(err)
	}
	for _, mt := range mainT.Templates() {
		tn := mt.Name()
		if strings.HasSuffix(tn, ".tmpl") {
			// skip
			continue
		}
		b := new(bytes.Buffer)
		err = mt.Execute(b, graph)
		if err != nil {
			return err
		}
		fname := filepath.Join(e.TargetDir, mt.Name()+".go")
		os.MkdirAll(filepath.Dir(fname), os.ModePerm)
		if err = gen.FormatGoFile(fname, b.Bytes()); err != nil {
			return err
		}
	}
	err = runGoModTidy(e.TargetDir)
	if err != nil {
		log.Println(err)
	}
	return nil
}

// webStaticAssets copy static assets to web/src/assets
func (e *Extension) webStaticAssets(ex gen.Extension) error {
	if !e.genWeb {
		return nil
	}
	// copy static assets
	return e.copyEmbedDir("template/web/src/assetsstatic", "web/src/assets")
}

func (e *Extension) copyEmbedDir(src, tar string) error {
	fl, err := templateDir.ReadDir(src)
	if err != nil {
		return err
	}
	for _, sf := range fl {
		if sf.IsDir() {
			if err = os.MkdirAll(filepath.Join(tar, sf.Name()), os.ModePerm); err != nil {
				return err
			}
			if err = e.copyEmbedDir(path.Join(src, sf.Name()), filepath.Join(tar, sf.Name())); err != nil {
				return err
			}
			continue
		}
		bts, err := templateDir.ReadFile(path.Join(src, sf.Name()))
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(tar, sf.Name()), bts, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func New(opt ...Option) *Extension {
	ex := &Extension{}
	ex.initTemplates()
	for _, o := range opt {
		o(ex)
	}
	return ex
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
		gen.ParseT("template/gomod.tmpl", templateDir, ExtensionFuncs),
	)
}

func runGoModTidy(dir string) error {
	err := gen.RunCmd(dir, "go", "mod", "tidy")
	if err != nil {
		return err
	}
	// wait for go mod tidy change go.mod and go.sum
	time.Sleep(1 * time.Second)
	return nil
}
