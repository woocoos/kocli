package schema

// Component 组件接口
type Component interface {
	// GetComponentName 获取组件名称
	GetComponentName() string
}

// LowCodeComponent 低代码编码UI组件.如 Page等属于编辑器容器
type LowCodeComponent struct {
	// 研发模式
	DevMode string `json:"devMode"`
	// 组件名称
	ComponentName string `json:"componentName"`
}

func (l LowCodeComponent) GetComponentName() string {
	return l.ComponentName
}

// ProCodeComponent 业务代码组件包信息.
type ProCodeComponent struct {
	NpmInfo
}

func (l ProCodeComponent) GetComponentName() string {
	return l.ComponentName
}
