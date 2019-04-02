package common

const (
	DEBUG = "debug"
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
	FATAL = "fatal"
)

const (
	ZERO = iota
	ONE
	TWO
	THREE
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
	NINE
	TEN
)

const (
	SUCCESS = iota + 200
	UNKNOWN
	ParamError
	UsernameNotExist
	UsernameExist
	PasswordError
	CodeError
	DataNotExist
	DataExist
	SystemBusy
	PermissionError
)

const (
	SuccessMessage          = "成功"
	UnknownMessage          = "未知错误"
	ParamErrorMessage       = "参数错误：[%v]"
	UsernameNotExistMessage = "用户名不存在"
	UsernameExistMessage    = "用户名已存在：[%v]"
	PasswordErrorMessage    = "用户名或密码错误"
	CodeErrorMessage        = "验证码错误"
	DataNotExistMessage     = "数据不存在：[%v]"
	DataExistMessage        = "数据已存在"
	SystemBusyMessage       = "服务器繁忙"
	PermissionErrorMessage  = "暂无权限"
)

const (
	TraceIdKey   = "traceId"
	CtxKey       = "commonContext"
	PasswordSalt = "passwordSalt"
)
