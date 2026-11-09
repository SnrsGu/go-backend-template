package middleware

import (
	"go-backend-template/internal/platform/apperror"
	"go-backend-template/internal/platform/response"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(
						"HTTP request panic",
						"method", r.Method,
						"url", r.URL.String(),
						"err", err,
						"stack", string(debug.Stack()),
					)
					err = response.Failure(w, apperror.New(apperror.CodeInternal, ""))
					if err != nil {
						logger.Error(
							"HTTP request failure",
							"error", err,
						)
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
