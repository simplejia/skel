package lib

const (
	// PROD 正式环境
	PROD = "prod"
	// DEV 开发环境，一般用于本机测试
	DEV = "dev"
	// TEST 测试环境
	TEST = "test"
)

// Resp 用于定义返回数据格式(json)
type Resp struct {
	Ret    Code        `json:"ret"`
	Msg    string      `json:"msg,omitempty"`
	Detail string      `json:"detail,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}
