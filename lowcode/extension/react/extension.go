package react

import (
	"github.com/tsingsun/woocoo/cmd/woco/gen"
)

var _ gen.Extension = (*Extension)(nil)

type Option func(*Extension)

type Extension struct {
	hooks     []gen.Hook
	templates []*gen.Template
}

func New(opt ...Option) *Extension {
	ex := &Extension{}
	ex.initTemplates()
	for _, o := range opt {
		o(ex)
	}
	return ex
}

func (e *Extension) GeneratedHooks() []gen.GeneratedHook {
	return []gen.GeneratedHook{
		func(ex gen.Extension) error {
			return nil
		},
	}
}

func (e *Extension) Name() string {
	return "react"
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
		gen.ParseT("template/container.tmpl", templateDir, ExtensionFuncs),
	)
}
