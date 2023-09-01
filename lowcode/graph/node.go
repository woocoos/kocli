package graph

import (
	"fmt"
	"github.com/woocoos/kocli/lowcode/schema"
	"sort"
	"strings"
)

// Node 相当于一个页面
type Node struct {
	*Config
	project *schema.ProjectSchema
	Page    *schema.PageSchema
}

// Import 导入包
type Import struct {
	Package      string            // 包名
	Default      string            // default export
	Component    map[string]string //组件名:导出名
	SubComponent map[string]string //子组件名:导出名
}

func (i *Import) String() string {
	var (
		val           []string
		destructuring string
	)
	for com, as := range i.Component {
		if as == com {
			val = append(val, com)
		} else {
			val = append(val, com+" as "+as)
		}
	}
	if len(val) > 0 {
		destructuring = fmt.Sprintf("{ %s }", strings.Join(val, ", "))
		if i.Default != "" {
			destructuring = ", " + destructuring
		}
	}
	return fmt.Sprintf(`import %s%s from '%s'`, i.Default, destructuring, i.Package)
}

func NewNode(c *Config, p *schema.ProjectSchema, page *schema.PageSchema) *Node {
	return &Node{
		Config:  c,
		project: p,
		Page:    page,
	}
}

// 添加NPM包
func (n *Node) addPkg(info schema.NpmInfo) {
	for _, pkg := range n.Packages {
		if pkg.Package == info.Package {
			return
		}
	}
	n.Packages = append(n.Packages, &schema.NpmInfo{
		Package: info.Package,
		Version: info.Version,
	})
}

// Imports 获取需要导入的包.目前默认直接从schema中获取
//
// 参考: 低代码引擎搭建协议: 组件映射关系
func (n *Node) Imports() []*Import {
	var imps = make(map[string]*Import)
	for _, com := range n.Page.LoadComponentNames() {
		npm := n.project.FindComponentMap(com)
		if npm == nil {
			continue
		}
		info, ok := imps[npm.Package]
		if !ok {
			info = &Import{
				Package:   npm.Package,
				Component: make(map[string]string),
			}
			imps[npm.Package] = info
			n.addPkg(npm.NpmInfo)
		}
		if npm.Destructuring {
			info.Component[npm.ExportName] = npm.ComponentName
		} else {
			info.Default = npm.ComponentName
		}
		if npm.SubName != "" {
			info.SubComponent[npm.ExportName+"."+npm.SubName] = npm.ComponentName
		}
	}
	if len(imps) == 0 {
		return nil
	}
	var imports = make([]*Import, 0, len(imps))
	for _, imp := range imps {
		imports = append(imports, imp)
	}
	// sort
	sort.Slice(imports, func(i, j int) bool {
		return imports[i].Package < imports[j].Package
	})
	return imports
}
