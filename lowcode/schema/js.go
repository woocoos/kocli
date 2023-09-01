package schema

// JSExpression 变量表达式
type JSExpression struct {
	Type string `json:"type"`
	// 表达式字符串
	Value string `json:"value"`
	// TODO 模似值
	Mock string `json:"mock,omitempty"`
	// TODO 源码
	Compiled string `json:"compiled,omitempty"`
}

// NodeDataType implement NodeData
func (jse *JSExpression) NodeDataType() string {
	return "JSExpression"
}

// JSFunction 事件函数类型
// @see https://lowcode-engine.cn/lowcode
//
// / 保留与原组件属性、生命周期 ( React / 小程序) 一致的输入参数，并给所有事件函数 binding 统一一致的上下文（当前组件所在容器结构的 this 对象）
type JSFunction struct {
	Type string `json:"type"`
	// 函数定义，或直接函数表达式
	Value string `json:"value"`
	// TODO 模似值
	Mock string `json:"mock,omitempty"`
	// TODO 源码
	Compiled string `json:"compiled,omitempty"`
	// 额外扩展属性，如 extType、events
	// @todo 待标准描述
	Others map[string]any `json:"-"`
}

// NodeDataType implement NodeData
func (jsf *JSFunction) NodeDataType() string {
	return "JSFunction"
}

// JSSlot Slot函数类型,通常用于描述组件的某一个属性为 ReactNode 或 Function return ReactNode 的场景
//
// 如果组件中Params为空,则Value当对象处理,否则当函数处理
type JSSlot struct {
	Type string `json:"type"`
	// todo 待标准描述
	ID string `json:"id,omitempty"`
	// todo 待标准描述
	Name string `json:"name,omitempty"`
	// todo 待标准描述
	Title string `json:"title,omitempty"`
	// 组件的某一个属性为 Function return ReactNode 时，函数的入参
	Params []string `json:"params,omitempty"`
	// 具体的值
	Value CompositeValue `json:"value,omitempty"`
}

func (jss *JSSlot) LoadComponentNames() []string {
	return jss.Value.LoadComponentNames()
}

// NodeDataType implement NodeData
func (jss *JSSlot) NodeDataType() string {
	return "JSSlot"
}
