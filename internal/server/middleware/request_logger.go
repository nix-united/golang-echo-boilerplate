package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type tracer interface {
	Start(ctx context.Context) (context.Context, error)
}

// requestLogger is a logging middleware that generated trace ID for each request.
type requestLogger struct {
	tracer tracer
}

func NewRequestLogger(tracer tracer) echo.MiddlewareFunc {
	return (&requestLogger{tracer: tracer}).handle
}

// handle creates trace and logs request information.
func (l *requestLogger) handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, err := l.tracer.Start(c.Request().Context())
		if err != nil {
			return fmt.Errorf("trace starter: %w", err)
		}

		c.SetRequest(c.Request().WithContext(ctx))

		errNext := next(c)

		level := slog.LevelInfo
		if c.Response().Status >= http.StatusInternalServerError {
			level = slog.LevelError
		}

		attrs := []any{
			"method", c.Request().Method,
			"status", c.Response().Status,
			"path", c.Path(),
		}

		if errNext != nil {
			attrs = append(attrs, "error", errNext.Error())
		}

		slog.Log(c.Request().Context(), level, "Request", slog.Group("http", attrs...))

		if errNext != nil {
			return fmt.Errorf("handle request with request logger: %w", errNext)
		}

		return nil
	}
}
