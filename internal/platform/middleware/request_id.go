package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
)

type requestIDKey struct{}

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var idBytes [16]byte
			if _, err := rand.Read(idBytes[:]); err != nil {
				panic("生成request_id失败: " + err.Error())
			}
			requestID := hex.EncodeToString(idBytes[:])
			w.Header().Set("X-Request-ID", requestID)
			ctx := context.WithValue(r.Context(), requestIDKey{}, requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetRequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDKey{}).(string)
	return requestID
}

func LoggerWithRequest(logger *slog.Logger, r *http.Request) *slog.Logger {
	return logger.With("request_id", GetRequestID(r.Context()))
}
