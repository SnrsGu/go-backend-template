package middleware

import (
	"go-backend-template/internal/platform/response"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recovery(logger *slog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(
						"HTTP请求panic",
						"method", r.Method,
						"path", r.URL.Path,
						"panic", err,
						"stack", string(debug.Stack()),
					)
					if err := response.Fail(w, response.CodeInternal); err != nil {
						logger.Error("错误响应失败", "error", err)
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
