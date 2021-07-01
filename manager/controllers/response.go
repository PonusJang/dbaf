package controllers

const (
	// common

	CODE_SUCCESS     = 1000
	CODE_FALURE      = 1001
	CODE_ERROR       = 1002
	CODE_ILLEGE      = 1003
	CODE_PARAM_ERROR = 1004

	CODE_LOGIN_SUCCESS        = 2000
	CODE_LOGIN_FAILURE        = 2001
	CODE_TOKEN_VERIFY_SECCESS = 2002
	CODE_TOKEN_VERIFY_FAILURE = 2003
)

const (
	MSG_SUCCESS         = "成功"
	MSG_FAILURE         = "失败"
	MSG_ERROR           = "错误"
	MSG_PARAM_ERROR     = "参数错误"
	MSG_USER_PASS_ERROR = "用户名密码错误"
	MSG_TOKEN_ERROR     = "Token错误"
	MSG_TOKEN_EXPIRE    = "Token超时"
)

// 统一响应状态

type Ret struct {
	Code   int         `json:code`
	Status bool        `json:status`
	Msg    string      `json:msg`
	Data   interface{} `json:data`
}

type ResData struct {
	Count int
	Data  interface{}
}
