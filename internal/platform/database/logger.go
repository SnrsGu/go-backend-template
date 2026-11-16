package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type gormLogger struct {
	logger        *slog.Logger
	level         gormlogger.LogLevel
	slowThreshold time.Duration
}

func newGormLogger(logger *slog.Logger, slowThreshold time.Duration) gormlogger.Interface {
	return &gormLogger{
		logger:        logger,
		level:         gormlogger.Warn,
		slowThreshold: slowThreshold,
	}
}

func (l *gormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	copied := *l
	copied.level = level
	return &copied
}

func (l *gormLogger) Info(ctx context.Context, message string, args ...interface{}) {
	if l.level >= gormlogger.Info {
		l.logger.InfoContext(ctx, fmt.Sprintf(message, args...))
	}
}

func (l *gormLogger) Warn(ctx context.Context, message string, args ...interface{}) {
	if l.level >= gormlogger.Warn {
		l.logger.WarnContext(ctx, fmt.Sprintf(message, args...))
	}
}

func (l *gormLogger) Error(ctx context.Context, message string, args ...interface{}) {
	if l.level >= gormlogger.Error {
		l.logger.ErrorContext(ctx, fmt.Sprintf(message, args...))
	}
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level == gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	if err != nil && l.level >= gormlogger.Error {
		l.logger.ErrorContext(ctx, "数据库 SQL 执行失败", "duration_ms", elapsed.Milliseconds(), "error", err)
		return
	}

	if elapsed > l.slowThreshold &&
		l.level >= gormlogger.Warn {
		l.logger.WarnContext(ctx, "慢SQL", "duration_ms", elapsed.Milliseconds())
		return
	}

	if l.level >= gormlogger.Info {
		_, rows := fc()

		l.logger.InfoContext(ctx, "数据库 SQL 执行完成", "duration_ms", elapsed.Milliseconds(), "rows", rows)
	}
}
