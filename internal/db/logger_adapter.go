package db

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm/logger"
)

var _ logger.Interface = (*LoggerAdapter)(nil)

type LoggerAdapter struct{}

func newLoggerAdapter() *LoggerAdapter {
	return &LoggerAdapter{}
}

func (a *LoggerAdapter) LogMode(logger.LogLevel) logger.Interface {
	return a
}

func (a *LoggerAdapter) Info(ctx context.Context, message string, args ...any) {
	slog.InfoContext(ctx, message, "args", args)
}

func (a *LoggerAdapter) Warn(ctx context.Context, message string, args ...any) {
	slog.WarnContext(ctx, message, "args", args)
}

func (a *LoggerAdapter) Error(ctx context.Context, message string, args ...any) {
	slog.ErrorContext(ctx, message, "args", args)
}

func (a *LoggerAdapter) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rowsAffected := fc()

	args := []any{
		"begin", begin.Format(time.DateTime),
		"sql", sql,
		"rows_affected", rowsAffected,
	}
	if err != nil {
		args = append(args, "err", err.Error())
	}

	slog.DebugContext(ctx, "Trace DB query execution", args...)
}
