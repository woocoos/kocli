package gentool

import (
	"embed"
	"github.com/tsingsun/woocoo/cmd/woco/gen"
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
)

func initTemplates() {
	templates = gen.ParseT("templates", templateDir, nil, "template/*.tmpl")
}
