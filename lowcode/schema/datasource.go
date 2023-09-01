package schema

type DataSource struct {
	DataHandler JSFunction `json:"dataHandler,omitempty"`
}

type DataSourceConfigOptions struct {
	// 请求地址
	Uri CompositeValue `json:"uri,omitempty"`
	Api CompositeValue `json:"api,omitempty"`
	// 请求参数
	Params CompositeValue `json:"params,omitempty"`
	// 请求方法
	Method CompositeValue `json:"method,omitempty"`
	// 是否支持跨域
	IsCors CompositeValue `json:"isCors,omitempty"`
	// 超时时长,单位 ms
	Timeout CompositeValue `json:"timeout,omitempty"`
	// 请求头信息
	Headers CompositeValue            `json:"headers,omitempty"`
	Others  map[string]CompositeValue `json:"-"`
}

type DataSourceConfig struct {
	// 	数据请求 ID 标识
	ID string `json:"id"`
	// 是否为初始数据,值为 true 时，将在组件初始化渲染时自动发送当前数据请求
	IsInit CompositeValue `json:"isInit,omitempty"`
	// 是否需要串行执行,值为 true 时，当前请求将被串行执行
	IsSync CompositeValue `json:"isSync,omitempty"`
	// 数据请求类型,支持四种类型：fetch/mtop/jsonp/custom
	Type string `json:"type,omitempty"`
	// 自定义扩展的外部请求处理器,仅 type='custom' 时生效
	RequestHandler JSFunction `json:"requestHandler,omitempty"`
	// request 成功后的回调函数
	DataHandler JSFunction `json:"dataHandler,omitempty"`
	// request 失败后的回调函数
	ErrorHandler JSFunction `json:"errorHandler,omitempty"`
	WillFetch    JSFunction `json:"willFetch,omitempty"`
	ShouldFetch  JSFunction `json:"shouldFetch,omitempty"`
	// 请求参数
	Options DataSourceConfigOptions   `json:"options,omitempty"`
	Others  map[string]CompositeValue `json:"-"`
}
