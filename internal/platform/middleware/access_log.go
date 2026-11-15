package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func AccessLog(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			writer := getResponseWriter(w)
			next.ServeHTTP(writer, r)
			duration := time.Since(start)
			logger.Info(
				"请求完成",
				"request_id", GetRequestID(r.Context()),
				"method", r.Method,
				"path", r.URL.Path,
				"status", writer.Status(),
				"duration_ms", duration.Milliseconds(),
				"bytes", writer.BytesWritten(),
				"remote_addr", r.RemoteAddr,
			)
		})
	}
}
