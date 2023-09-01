package schema

import "encoding/json"

// NodeSchema 搭建基础协议 - 单个组件树节点描述
type NodeSchema struct {
	ID string `json:"id,omitempty"`
	// 组件名称 必填、首字母大写
	ComponentName string `json:"componentName"`
	// 组件属性对象
	Props *Props `json:"props"`
	// 渲染条件
	Condition CompositeValue `json:"condition,omitempty"`
	// 循环数据
	Loop CompositeValue `json:"loop,omitempty"`
	// 循环迭代对象、索引名称 ["item", "index"]
	LoopArgs [2]string `json:"loopArgs,omitempty"`
	// 子节点
	Children []*NodeData `json:"children,omitempty"`
	IsLocked bool        `json:"isLocked,omitempty"`

	// ------- future support -----
	ConditionGroup string `json:"conditionGroup,omitempty"`
	Title          string `json:"title,omitempty"`
	Ignore         bool   `json:"ignore,omitempty"`
	Locked         bool   `json:"locked,omitempty"`
	Hidden         bool   `json:"hidden,omitempty"`
	IsTopField     bool   `json:"isTopField,omitempty"`
}

func (ns *NodeSchema) GetComponentName() string {
	return ns.ComponentName
}

// LoadComponentNames 读取容器内的所有组件名称,根据该名称可以获取组件的包信息
func (ns *NodeSchema) LoadComponentNames() (cns []string) {
	if ns.ComponentName != "" {
		cns = []string{ns.ComponentName}
	}
	if ns.Props != nil {
		cns = append(cns, ns.Props.PropsMap.LoadComponentNames()...)
	}
	for _, child := range ns.Children {
		switch cv := child.Value.(type) {
		case *NodeSchema:
			cns = append(cns, cv.LoadComponentNames()...)
		}
	}
	return
}

func (ns *NodeSchema) UnmarshalJSON(bytes []byte) error {
	type ins NodeSchema
	var tmp ins
	if err := json.Unmarshal(bytes, &tmp); err != nil {
		return err
	}
	*ns = NodeSchema(tmp)
	if ns.Props != nil {
		for _, child := range ns.Props.Children {
			ns.Children = append(ns.Children, child)
		}
	}
	return nil
}

// Props 组件属性对象
type Props struct {
	// 组件 ID
	ID string `json:"id,omitempty"`
	// 组件样式类名
	ClassName CompositeValue    `json:"className,omitempty"`
	Style     CompositeValueMap `json:"style,omitempty"`
	// 组件 ref 名称,可通过 this.$(ref) 获取组件实例
	Ref      string            `json:"ref,omitempty"`
	PropsMap CompositeValueMap `json:"-"`
	// 在props中配置的children,独立存储
	Children []*NodeData `json:"-"`
}

func (p *Props) UnmarshalJSON(bytes []byte) error {
	var propsMap CompositeValueMap
	if err := json.Unmarshal(bytes, &propsMap); err != nil {
		return err
	}
	if id, ok := propsMap["id"]; ok {
		p.ID = id.Value.(string)
		delete(propsMap, "id")
	}
	if className, ok := propsMap["className"]; ok {
		p.ClassName = className
		delete(propsMap, "className")
	}
	if style, ok := propsMap["style"]; ok {
		p.Style = style.Value.(CompositeValueMap)
		delete(propsMap, "style")
	}
	if ref, ok := propsMap["ref"]; ok {
		p.Ref = ref.Value.(string)
		delete(propsMap, "ref")
	}
	if children, ok := propsMap["children"]; ok {
		p.Children = append(p.Children, &NodeData{Value: children.Value})
		delete(propsMap, "children")
	}
	p.PropsMap = propsMap
	return nil
}

func (p *Props) Prop(key string) string {
	if p.PropsMap == nil {
		return ""
	}
	if v, ok := p.PropsMap[key]; ok {
		return v.Value.(string)
	}
	return ""
}

// NodeData 节点数据类型
//
// 包括: NodeSchema, DOMText, JSExpression, I18NData及其数组类型
type NodeData CompositeValue

func (n *NodeData) UnmarshalJSON(bytes []byte) error {
	var c CompositeValue
	if err := json.Unmarshal(bytes, &c); err != nil {
		return err
	}
	*n = NodeData(c)
	return nil
}
