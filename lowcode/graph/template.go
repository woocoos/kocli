package graph

import (
	"embed"
	"github.com/Masterminds/sprig/v3"
	"github.com/tsingsun/woocoo/cmd/woco/gen"
	"path/filepath"
)

var (
	// templates hold the Go templates for the code generation.
	templates *gen.Template
	//go:embed template/*
	templateDir embed.FS

	PageTemplates = []PageTemplate{
		{
			Name: "index",
			Format: func(n *Node) string {
				pid := n.Page.FileName
				dir := n.Config.PageRoot()
				if IsValidNodeName(pid) {
					return filepath.Join(dir, pid, "index.tsx")
				}
				return filepath.Join(dir, "index.tsx")
			},
		},
	}
	GraphTemplates = []GraphTemplate{}
)

type PageTemplate struct {
	Name   string
	Format func(*Node) string
}

// GraphTemplate 对应每一个输出文件.
type GraphTemplate struct {
	Name   string
	Format func(*Graph) string
}

func initTemplates() {
	templates = gen.NewTemplate("templates")
	//templates.Delims("[[", "]]")
	gen.MustParse(templates.Funcs(Funcs).Funcs(sprig.FuncMap()).ParseFS(templateDir, "template/*/*.tmpl"))
}

// IsValidNodeName is true iff the node name matches the pattern of LabelNameRE. This
// method, however, does not use LabelNameRE for the check but a much faster
// hardcoded implementation.
func IsValidNodeName(name string) bool {
	if len(name) == 0 {
		return false
	}
	for i, b := range name {
		if !((b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_' || (b >= '0' && b <= '9' && i > 0)) {
			return false
		}
	}
	return true
}

func printCompositeValue() {

}
