package graph

import (
	"bytes"
	"fmt"
	"github.com/tsingsun/woocoo/cmd/woco/gen"
	"github.com/woocoos/kocli/lowcode/schema"
	"path/filepath"
	"text/template/parse"
)

type Config struct {
	// ProjectSchema文件路径
	Schema string `json:"schema"`
	// 目标目录, 默认为schema所在目录.对于knockout项目,请指定为`{{prjdir}}/web`
	Target    string `json:"target,omitempty"`
	Templates []*gen.Template
	Header    string `json:"header,omitempty"`
	// Hooks hold an optional list of Hooks to apply on the graph before/after the code-generation.
	Hooks          []gen.Hook
	GeneratedHooks []gen.GeneratedHook
	// 依赖包
	Packages []*schema.NpmInfo `json:"-"`
	// 依赖包 和node类似, 包名:版本号
	Dependencies map[string]string `json:"dependencies"`
}

// PageRoot 页面生成目录根据JS项目为src/page
func (c Config) PageRoot() string {
	return "src/page"
}

// Graph is a graph of code generation.
type Graph struct {
	*Config

	// lowcode schema
	ProjectSchema *schema.ProjectSchema
	// 页面组件.为顶级组件
	Pages []*Node
}

func NewGraph(c *Config, sch *schema.ProjectSchema) (g *Graph, err error) {
	g = &Graph{
		Config:        c,
		ProjectSchema: sch,
	}
	if g.Dependencies == nil {
		g.Dependencies = make(map[string]string)
	}
	for _, rootSchema := range g.ProjectSchema.ComponentsTree {
		g.Pages = append(g.Pages, NewNode(c, sch, rootSchema.(*schema.PageSchema)))
	}
	return g, nil
}

// Name implement gen.Graph interface
func (*Graph) Name() string {
	return "ProjectInit"
}

// Gen generates the artifacts for the graph.
func (g *Graph) Gen() error {
	return gen.ExecGen(generate, g)
}

func (g *Graph) Hooks() []gen.Hook {
	return nil
}

func (g *Graph) GeneratedHooks() []gen.GeneratedHook {
	return g.Config.GeneratedHooks
}

func (g *Graph) Templates() []*gen.Template {
	return g.Config.Templates
}

func (g *Graph) templates() (*gen.Template, []GraphTemplate) {
	initTemplates()
	var (
		roots = make(map[string]struct{})
	)
	gt := make([]GraphTemplate, 0, len(g.Config.Templates))
	for _, rootT := range g.Config.Templates {
		templates.Funcs(rootT.FuncMap)
		for _, tpl := range rootT.Templates() {
			if parse.IsEmptyTree(tpl.Root) {
				continue
			}
			name := tpl.Name()
			switch {
			case templates.Lookup(name) == nil:
				gt = append(gt, GraphTemplate{
					Name: name,
					Format: func(graph *Graph) string {
						return name
					},
				})
				roots[name] = struct{}{}
			}
			templates = gen.MustParse(templates.AddParseTree(name, tpl.Tree))
		}
	}
	return templates, gt
}

// 依赖信息
func (g *Graph) genDependencies() error {
	for _, b := range g.ProjectSchema.ComponentsMap {
		pc, ok := b.(*schema.ProCodeComponent)
		if !ok {
			continue
		}
		if _, ok := g.Dependencies[pc.Package]; !ok {
			g.Dependencies[pc.Package] = pc.Version
		}
	}
	return nil
}

// Actually execute the generated code.
func generate(gg gen.Extension) error {
	g := gg.(*Graph)
	var (
		assets   gen.Assets
		external []GraphTemplate
	)
	if err := g.genDependencies(); err != nil {
		return err
	}

	templates, external = g.templates()
	//pkg := g.Package
	assets.AddDir(filepath.Join(g.Config.Target))
	for _, page := range g.Pages {
		for _, tmpl := range PageTemplates {
			b := bytes.NewBuffer(nil)
			if dir := filepath.Dir(tmpl.Format(page)); dir != "." {
				assets.AddDir(filepath.Join(g.Config.Target, dir))
			}
			if err := templates.ExecuteTemplate(b, tmpl.Name, page); err != nil {
				return fmt.Errorf("execute template %q: %w", tmpl.Name, err)
			}
			assets.Add(filepath.Join(g.Target, tmpl.Format(page)), b.Bytes())
		}
	}
	for _, tmpl := range append(GraphTemplates, external...) {
		b := bytes.NewBuffer(nil)
		if dir := filepath.Dir(tmpl.Format(g)); dir != "." {
			assets.AddDir(filepath.Join(g.Config.Target, dir))
		}
		if err := templates.ExecuteTemplate(b, tmpl.Name, g); err != nil {
			return fmt.Errorf("execute template %q: %w", tmpl.Name, err)
		}
		assets.Add(filepath.Join(g.Target, tmpl.Format(g)), b.Bytes())
	}
	// Write and Format Assets only if template execution
	// finished successfully.
	if err := assets.Write(); err != nil {
		return err
	}
	assets.AddFormatter(func(path string, content []byte) error {
		return gen.RunCmd(g.Target, "tsfmt", "-r", path)
	}, ".tsx", ".ts")
	return assets.Format()
}
