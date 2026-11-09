package logger

import (
	"io"
	"log/slog"
)

func New(output io.Writer, level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: level, AddSource: true,
	})
	return slog.New(handler)
}
