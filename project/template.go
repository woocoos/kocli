package project

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
	templates = gen.ParseT("templates", templateDir, nil, "template/*.tmpl")
	tpkgs = gen.InitTemplates(templates, "import-knockout", Extension{})
	for k, v := range tpkgs {
		importPkg[k] = v
	}
}
