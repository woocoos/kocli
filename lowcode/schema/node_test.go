package schema

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/woocoos/kocli/test"
	"os"
	"testing"
)

func TestProps_UnmarshalJSON(t *testing.T) {
	type args struct {
		bytes string
	}
	tests := []struct {
		name    string
		args    args
		check   func(props Props)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "props",
			args: args{
				bytes: `
{
	"props": {
        "className": "btn",
        "style": {
        "width": 100,
        "height": 20
        },
        "text": "submit",
        "onClick": {
            "type": "JSFunction",
            "value": "function(e){console.log('hello world')}"
        },
		"title": {
			"type": "JSSlot",
			"value": [{
				"componentName": "Icon",
				"props": {}
			},{
				"componentName": "Icon",
				"props": {}
			}]
		}
    }
}`,
			},
			check: func(props Props) {
				assert.Equal(t, "btn", props.ClassName.Value.(string))
				assert.Equal(t, float64(100), props.Style["width"].Value.(float64))
				assert.Equal(t, float64(20), props.Style["height"].Value.(float64))
				assert.Equal(t, "submit", props.PropsMap["text"].Value.(string))
				assert.Equal(t, "function(e){console.log('hello world')}", props.PropsMap["onClick"].Value.(*JSFunction).Value)
				assert.Equal(t, "JSSlot", props.PropsMap["title"].Value.(*JSSlot).NodeDataType())
			},
			wantErr: assert.NoError,
		},
		{
			name: "props-children",
			args: args{
				bytes: `
{
	"props": {
		"columns": [
			{
				"title": "操作",
				"render": {
					"type": "JSSlot",
					"params": [
                        "text",
                        "record",
                        "index"
                    ],
					"value": [
						{
							"componentName": "Button"
						}
					]
				}
			}
		]
	}
}
`,
			},
			check: func(props Props) {
				col1 := props.PropsMap["columns"].Value.(CompositeValueSlice)[0].Value.(map[string]CompositeValue)
				assert.Equal(t, "JSSlot", col1["render"].Value.(*JSSlot).NodeDataType())
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			props := struct {
				Props Props `json:"props"`
			}{}
			err := json.Unmarshal([]byte(tt.args.bytes), &props)
			if !tt.wantErr(t, err) {
				return
			}
			tt.check(props.Props)
		})
	}
}

func TestNodeSchema_LoadComponentNames(t *testing.T) {
	type fields struct {
		ID            string
		ComponentName string
		Children      []*NodeData
		Props         *Props
	}
	tests := []struct {
		name    string
		fields  fields
		wantCns []string
	}{
		{
			name: "node",
			fields: fields{
				ComponentName: "Page",
			},
			wantCns: []string{"Page"},
		},
		{
			name: "node-children",
			fields: fields{
				ComponentName: "Page",
				Children: []*NodeData{
					{
						Value: CompositeValueSlice{
							{
								Value: &NodeSchema{
									ComponentName: "Button",
								},
							},
							{
								Value: &NodeSchema{
									ComponentName: "Input",
								},
							},
						},
					},
				},
			},
			wantCns: []string{"Page", "Button", "Input"},
		},
		{
			name: "from parse",
			fields: func() fields {
				bs, err := os.ReadFile(test.Path("testdata/protable/schema.json"))
				require.NoError(t, err)
				ps, err := ParseProjectSchema(bs)
				return fields{
					ComponentName: ps.ComponentsTree[0].GetRootSchemaName(),
					Children:      ps.ComponentsTree[0].(*PageSchema).Children,
					Props:         ps.ComponentsTree[0].(*PageSchema).Props,
				}
			}(),
			wantCns: []string{"Page", "ProTable", "Button", "ProPopconfirm", "Icon"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := &NodeSchema{
				ID:            tt.fields.ID,
				ComponentName: tt.fields.ComponentName,
				Children:      tt.fields.Children,
				Props:         tt.fields.Props,
			}
			assert.Equalf(t, tt.wantCns, ns.LoadComponentNames(), "LoadComponentNames()")
		})
	}
}
