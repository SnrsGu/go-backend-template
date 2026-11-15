package response

import (
	"encoding/json/v2"
	"fmt"
	"net/http"
)

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Success(w http.ResponseWriter, data any) error {
	return writeJSON(
		w,
		CodeSuccess.HTTPStatus(),
		Result{
			Code:    CodeSuccess.Value(),
			Message: CodeSuccess.Message(),
			Data:    data,
		},
	)
}

func Fail(w http.ResponseWriter, err error) error {
	code := ResolveCode(err)
	return writeJSON(
		w,
		code.HTTPStatus(),
		Result{
			Code:    code.Value(),
			Message: code.Message(),
			Data:    nil,
		},
	)
}

func writeJSON(w http.ResponseWriter, status int, result Result) error {
	body, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("响应序列化失败: %w", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(body); err != nil {
		return fmt.Errorf("写入响应失败: %w", err)
	}
	return nil
}
