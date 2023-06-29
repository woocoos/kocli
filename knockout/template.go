package knockout

import (
	"embed"
	"github.com/tsingsun/woocoo/cmd/woco/gen"
	"text/template"
)

var (
	// templates holds the Go templates for the code generation.
	templates *gen.Template
	//go:embed template/*
	templateDir embed.FS
	importPkg   = map[string]string{
		"context": "context",
		"errors":  "errors",
		"fmt":     "fmt",
		"math":    "math",
		"strings": "strings",
		"time":    "time",
		"regexp":  "regexp",
	}
	ExtensionFuncs = template.FuncMap{}
)

func initTemplates() {
	tpkgs := make(map[string]string)
	templates, tpkgs = gen.InitTemplates(templateDir, "import-knockout",
		Extension{},
		"template/*.tmpl")
	for k, v := range tpkgs {
		importPkg[k] = v
	}
}
