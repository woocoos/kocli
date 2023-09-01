package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseProjectSchema(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
		check   func(*ProjectSchema)
	}{
		{
			name: "componentsMap",
			args: args{
				str: `
{
	"version": "1.0.0",
	"componentsMap": [
		{
			"devMode": "lowcode",
			"componentName": "lowcode"
		},
		{
			"componentName": "procode",
			"package": "@tencent/lowcode",
			"version": "1.0.0",
			"destructuring": true,
			"exportName": "default",
			"subName": "sub",
			"main": "index.js"
		}
	]
}
`,
			},
			wantErr: assert.NoError,
			check: func(s *ProjectSchema) {
				assert.Len(t, s.ComponentsMap, 2)
				assert.Equal(t, "lowcode", s.ComponentsMap[0].GetComponentName())
				assert.Equal(t, "procode", s.ComponentsMap[1].GetComponentName())
			},
		},
		{
			name: "componentsTree",
			args: args{
				str: `
{
	"version": "1.0.0",
	"componentsTree": [{
		"componentName": "Page",
		"fileName": "page-1",
		"props": {
			"className": "page",
			"style": {
				"width": 100,
				"height": 100
			}
		},
		"children": [{
            "componentName": "Button",
            "props": {
                "text": {
					"type": "JSExpression",
                    "value": "this.state.btnText"
            	}
			}
        }],
		"state": {
			"btnText": "hello world"
		},
		"css": "body {font-size: 12px;}"
	},{
		"componentName": "Block",
		"fileName": "block-1",
		"props": {}
	}]
}
`,
			},
			wantErr: assert.NoError,
			check: func(s *ProjectSchema) {
				assert.Len(t, s.ComponentsTree, 2)
				assert.Equal(t, RootSchemaNamePage, s.ComponentsTree[0].GetRootSchemaName())
				assert.Equal(t, RootSchemaNameBlock, s.ComponentsTree[1].GetRootSchemaName())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseProjectSchema([]byte(tt.args.str))
			if !tt.wantErr(t, err) {
				return
			}
			if tt.check != nil {
				tt.check(got)
			}
		})
	}
}
