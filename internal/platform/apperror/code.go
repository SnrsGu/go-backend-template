package apperror

import "net/http"

type Code int

const (
	CodeSuccess         Code = 0
	CodeInternal        Code = 10000
	CodeInvalidArgument Code = 10001
	CodeNotFound        Code = 10002

	//CodeUserNotFound Code = 20001
)

type CodeDefinition struct {
	HTTPStatus int
	Message    string
}

var codeDefinitions = map[Code]CodeDefinition{
	CodeSuccess:         {HTTPStatus: http.StatusOK, Message: "success"},
	CodeInternal:        {HTTPStatus: http.StatusInternalServerError, Message: "未知错误"},
	CodeInvalidArgument: {HTTPStatus: http.StatusBadRequest, Message: "参数错误"},
	CodeNotFound:        {HTTPStatus: http.StatusNotFound, Message: "数据不存在"},
}

func (c Code) Definition() CodeDefinition {
	definition, exists := codeDefinitions[c]
	if !exists {
		definition = codeDefinitions[CodeInternal]
	}
	if definition.HTTPStatus == 0 {
		definition.HTTPStatus = http.StatusOK
	}
	return definition
}

func (c Code) HTTPStatus() int {
	status := c.Definition().HTTPStatus
	if status == 0 {
		return http.StatusOK
	}
	return status
}

func (c Code) Message() string {
	return c.Definition().Message
}

func (c Code) Normalize() Code {
	if _, exists := codeDefinitions[c]; !exists {
		return CodeInternal
	}
	return c
}
