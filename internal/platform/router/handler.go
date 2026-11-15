package router

import (
	"go-backend-template/internal/platform/response"
	"log/slog"
	"net/http"
)

type HandlerFunc func(writer http.ResponseWriter, request *http.Request) error

func wrapHandler(logger *slog.Logger, handler HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := handler(writer, request)
		if err == nil {
			return
		}
		code := response.ResolveCode(err)
		logError(logger, request, code, err)
		if writeErr := response.Fail(writer, err); writeErr != nil {
			logger.Error("错误响应失败", "error", writeErr)
		}
	}
}

func logError(logger *slog.Logger, request *http.Request, code response.Code, err error) {
	attributes := []any{
		"method", request.Method,
		"path", request.URL.Path,
		"code", code.Value(),
		"http_status", code.HTTPStatus(),
		"error", err,
	}

	if code.HTTPStatus() >= http.StatusInternalServerError {
		logger.Error("请求处理失败", attributes...)
		return
	}

	logger.Warn("HTTP 请求失败", attributes...)
}
