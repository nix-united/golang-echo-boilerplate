package slogx

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"

	"github.com/google/uuid"
)

type trace struct {
	trace string
	index *atomic.Int64
}

type TraceStarter struct {
	newUUID func() (uuid.UUID, error)
}

func NewTraceStarter(newUUID func() (uuid.UUID, error)) *TraceStarter {
	return &TraceStarter{newUUID: newUUID}
}

func (s *TraceStarter) Start(ctx context.Context) (context.Context, error) {
	traceID, err := s.newUUID()
	if err != nil {
		return nil, fmt.Errorf("new trace id: %w", err)
	}

	return withTrace(ctx, trace{trace: traceID.String(), index: &atomic.Int64{}}), nil
}

var _ slog.Handler = (*traceHandler)(nil)

type traceHandler struct {
	handler slog.Handler
}

func newTraceHandler(handler slog.Handler) *traceHandler {
	return &traceHandler{handler: handler}
}

func (h *traceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *traceHandler) Handle(ctx context.Context, record slog.Record) error {
	record = h.addAttrs(ctx, record)

	if err := h.handler.Handle(ctx, record); err != nil {
		return fmt.Errorf("handle log with trace: %w", err)
	}

	return nil
}

func (h *traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newTraceHandler(h.handler.WithAttrs(attrs))
}

func (h *traceHandler) WithGroup(name string) slog.Handler {
	return newTraceHandler(h.handler.WithGroup(name))
}

func (h *traceHandler) addAttrs(ctx context.Context, record slog.Record) slog.Record {
	if ctx == nil {
		return record
	}

	t, ok := traceFromContext(ctx)
	if !ok {
		return record
	}

	record.AddAttrs(slog.Group("trace", slog.String("trace", t.trace), slog.Int64("index", t.index.Add(1))))

	return record
}

type traceKeyType int8

var traceKey traceKeyType = 1

func withTrace(ctx context.Context, trace trace) context.Context {
	return context.WithValue(ctx, traceKey, trace)
}

func traceFromContext(ctx context.Context) (trace, bool) {
	trace, ok := ctx.Value(traceKey).(trace)
	return trace, ok
}
