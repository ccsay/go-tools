package result

import "github.com/liuchonglin/go-tools/common"

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
}

func NewSuccess() *Result {
	return &Result{Code: common.SUCCESS, Message: common.SuccessMessage}
}

func NewError(code int, message string) *Result {
	return &Result{Code: code, Message: message}
}

func NewSuccessData(data interface{}) *Result {
	return &Result{Code: common.SUCCESS, Message: common.SuccessMessage, Data: data}
}

func NewSuccessPage(data interface{}, total int64) *Result {
	return &Result{Code: common.SUCCESS, Message: common.SuccessMessage, Data: data, Total: total}
}
