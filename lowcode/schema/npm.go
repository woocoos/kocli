package schema

// NpmInfo 描述组件映射关系的集合
type NpmInfo struct {
	//源码组件名称
	ComponentName string `json:"componentName,omitempty"`
	//源码组件库名
	Package string `json:"package"`
	//源码组件版本号
	Version string `json:"version,omitempty"`
	// 使用解构方式对模块进行导出
	Destructuring bool `json:"destructuring,omitempty"`
	// 包导出的组件名
	ExportName string `json:"exportName,omitempty"`
	// 下标子组件名称
	SubName string `json:"subName,omitempty"`
	// 包导出组件入口文件路径
	Main string `json:"main,omitempty"`
}
