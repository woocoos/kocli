package knockout

import (
	"github.com/tsingsun/woocoo/cmd/woco/gen"
)

var _ gen.Extension = (*Extension)(nil)

type Option func(*Extension)

func WithTargetDir(dir string) Option {
	return func(e *Extension) {
		e.TargetDir = dir
	}
}

type Extension struct {
	TargetDir string
	hooks     []gen.Hook
	templates []*gen.Template
}

func (e *Extension) GeneratedHooks() []gen.GeneratedHook {
	return []gen.GeneratedHook{
		func(ex gen.Extension) error {
			// return helper.RunCmd(e.TargetDir, "go", "run", "codegen/entgen/entc.go")
			return nil
		},
	}
}

func New(opt ...Option) *Extension {
	ex := &Extension{}
	for _, o := range opt {
		o(ex)
	}
	ex.initTemplates()
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
		gen.ParseT("template/app.yaml.tmpl", templateDir, ExtensionFuncs),
		gen.ParseT("template/codegen.tmpl", templateDir, ExtensionFuncs),
		gen.ParseT("template/makefile.tmpl", templateDir, ExtensionFuncs),
		gen.ParseT("template/test.tmpl", templateDir, ExtensionFuncs),
		gen.ParseT("template/main.tmpl", templateDir, ExtensionFuncs),
		gen.ParseT("template/gomod.tmpl", templateDir, ExtensionFuncs),
	)
}
