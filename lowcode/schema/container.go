package schema

var (
	_ RootSchema = (*PageSchema)(nil)
	_ RootSchema = (*BlockSchema)(nil)
	_ RootSchema = (*ComponentSchema)(nil)
)

var (
	_ Container = (*NodeSchema)(nil)
	_ Container = (*JSSlot)(nil)
	_ Container = (*CompositeValueSlice)(nil)
	_ Container = (*CompositeValueMap)(nil)
)

// Container 是具有放置子组件能力的组件.如NodeSchema,JSSlot等.
type Container interface {
	LoadComponentNames() []string
}

// PageSchema 页面,Page是一般做为业务组件的根存在
type PageSchema struct {
	ContainerSchema `json:",inline"`

	project *ProjectSchema
}

func NewPageSchema(schema *ProjectSchema) *PageSchema {
	return &PageSchema{
		project: schema,
	}
}

func (p PageSchema) GetID() string {
	return p.ID
}

func (p PageSchema) GetRootSchemaName() string {
	return RootSchemaNamePage
}

// BlockSchema 区块容器
type BlockSchema struct {
	ContainerSchema `json:",inline"`
}

func (b BlockSchema) GetID() string {
	return b.ID
}

func (b BlockSchema) GetRootSchemaName() string {
	return RootSchemaNameBlock
}

// ComponentSchema 低代码业务组件容器
type ComponentSchema struct {
	ContainerSchema `json:",inline"`
}

func (c ComponentSchema) GetID() string {
	return c.ID
}

func (c ComponentSchema) GetRootSchemaName() string {
	return RootSchemaNameComponent
}

// ContainerSchema 容器结构描述
//
// ComponentName: 'Block' | 'Page' | 'Component'
type ContainerSchema struct {
	NodeSchema `json:",inline"`
	// 文件名称
	FileName string         `json:"fileName"`
	Meta     map[string]any `json:"meta,omitempty"`
	// 容器初始数据
	State map[string]CompositeValue `json:"state,omitempty"`
	// 自定义方法设置,CompositeValue = JSExpression | JSFunction
	Methods map[string]CompositeValue `json:"methods,omitempty"`
	// 样式文件
	Css string `json:"css,omitempty"`
	// 异步数据源配置
	DataSource *DataSource `json:"dataSource,omitempty"`
	// 低代码业务组件默认属性
	DefaultProps map[string]CompositeValue `json:"defaultProps,omitempty"`
}
