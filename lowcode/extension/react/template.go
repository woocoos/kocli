package react

import (
	"embed"
	"github.com/tsingsun/woocoo/cmd/woco/gen"
	"text/template"
)

var (
	// templates hold the Go templates for the code generation.
	templates *gen.Template
	//go:embed template/*
	templateDir embed.FS

	ExtensionFuncs = template.FuncMap{}
)

func initTemplates() {
	templates = gen.NewTemplate("templates")
	//templates.Delims("[[", "]]")
	gen.MustParse(templates.Funcs(ExtensionFuncs).ParseFS(templateDir, "template/*.tmpl"))
}
