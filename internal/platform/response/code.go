package response

import (
	"errors"
	"fmt"
	"net/http"
)

type Code struct {
	value      int
	message    string
	httpStatus int
}

func NewCode(value int, message string, httpStatus ...int) Code {
	if value < 0 {
		panic("错误码不可小于0")
	}
	if message == "" {
		panic("错误信息不能为空")
	}
	if len(httpStatus) > 1 {
		panic("一个错误码只能绑定一个HTTP状态码")
	}
	status := http.StatusOK
	if len(httpStatus) == 1 {
		status = httpStatus[0]
	}
	if status < 200 || status > 599 ||
		status == http.StatusNoContent ||
		status == http.StatusResetContent ||
		status == http.StatusNotModified {
		panic(fmt.Sprintf("HTTP状态码不适用与JSON响应: %d", status))
	}
	return Code{
		value:      value,
		message:    message,
		httpStatus: status,
	}
}

func (c Code) Error() string {
	return c.message
}

func (c Code) Value() int {
	return c.value
}

func (c Code) Message() string {
	return c.message
}

func (c Code) HTTPStatus() int {
	return c.httpStatus
}

func ResolveCode(err error) Code {
	if err == nil {
		return CodeSuccess
	}
	if code, ok := errors.AsType[Code](err); ok {
		return code
	}
	return CodeInternal
}

var (
	CodeSuccess         = NewCode(0, "success")
	CodeInternal        = NewCode(10000, "未知错误", http.StatusInternalServerError)
	CodeInvalidArgument = NewCode(10001, "请求参数错误", http.StatusBadRequest)
	CodeUnauthorized    = NewCode(10002, "未登录或登录已失效", http.StatusUnauthorized)
	CodeForbidden       = NewCode(10003, "无权限访问", http.StatusForbidden)
)
