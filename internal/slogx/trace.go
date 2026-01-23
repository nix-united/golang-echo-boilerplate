package slogx

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"

	"github.com/google/uuid"
)

type trace struct {
	// traceID represents a common UUID that share logs during request processing.
	traceID string

	// spanID represents a number of a log within a span.
	spanID *atomic.Int64

	// baggage represents an additional information that log messages shares.
	// For example, it could be a user ID.
	baggage map[string]any
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

	return contextWithTrace(ctx, trace{
		traceID: traceID.String(),
		spanID:  &atomic.Int64{},
		baggage: make(map[string]any),
	}), nil
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

	attrs := []any{
		slog.String("trace_id", t.traceID),
		slog.Int64("span_id", t.spanID.Add(1)),
	}

	for key, value := range t.baggage {
		attrs = append(attrs, slog.Any(key, value))
	}

	record.AddAttrs(slog.Group("trace", attrs...))

	return record
}

type traceKey struct{}

func contextWithTrace(ctx context.Context, trace trace) context.Context {
	return context.WithValue(ctx, traceKey{}, trace)
}

func traceFromContext(ctx context.Context) (trace, bool) {
	trace, ok := ctx.Value(traceKey{}).(trace)
	return trace, ok
}

// ContextWithBaggage appends [key] field with [value] to all log messages.
func ContextWithBaggage(ctx context.Context, key string, value any) context.Context {
	t, ok := traceFromContext(ctx)
	if !ok {
		return ctx
	}
	t.baggage[key] = value
	return contextWithTrace(ctx, t)
}

// ContextWithUserID appends user_id field to all log messages.
func ContextWithUserID(ctx context.Context, userID uint) context.Context {
	return ContextWithBaggage(ctx, "user_id", userID)
}
