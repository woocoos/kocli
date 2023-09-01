package graph

import (
	"fmt"
	"github.com/tsingsun/woocoo/cmd/woco/gen"
	"github.com/woocoos/kocli/lowcode/schema"
	"text/template"
)

var Funcs = template.FuncMap{
	"hasKey":         hasKey,
	"isComponent":    isComponent,
	"propsNeedBrace": propsNeedBrace,
	"sQuoteToQuote":  sQuoteToQuote,
	"joinSlice":      gen.Funcs["join"],
	"debugger":       debugger,
}

// debugger use {{ debugger . }} in template
func debugger(t any) string {
	fmt.Print(t)
	return ""
}

func hasKey(m map[string]any, key string) bool {
	_, ok := m[key]
	return ok
}

func isComponent(tp any) bool {
	_, ok := tp.(schema.Component)
	return ok
}

// 判断组件Props是否需要添加大括号,对于字符串类型不需要
func propsNeedBrace(v schema.CompositeValue) bool {
	str := v.TypeString()
	switch schema.CVType(str) {
	case schema.CVTypeString:
		return false
	}
	return true
}

// 单引号转双引号,在组件属性为字符串时,按约定使用双引号.
func sQuoteToQuote(s string) string {
	if s[0] != '\'' {
		return s
	}
	return `"` + s[1:len(s)-1] + `"`
}
