// Package logger собирает *slog.Logger приложения и подмешивает в каждую
// запись request_id из context.Context.
//
// Логгер создаётся один раз в main и передаётся вниз через зависимости.
// request_id кладёт в контекст HTTP-middleware, а логгер достаёт его сам,
// поэтому в хендлерах достаточно писать через *Context-методы:
//
//	h.log.ErrorContext(r.Context(), "db ping failed", "err", err)
package logger

import (
	"context"
	"io"
	"log/slog"
)

type requestIDKey struct{}

// WithRequestID возвращает контекст с request_id.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// RequestID достаёт request_id из контекста; ok == false, если его нет.
func RequestID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey{}).(string)
	return id, ok && id != ""
}

// New создаёт JSON-логгер, который пишет в w и добавляет request_id из контекста.
func New(w io.Writer, level slog.Leveler) *slog.Logger {
	return slog.New(NewContextHandler(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level})))
}

// NewContextHandler оборачивает любой slog.Handler: перед записью добавляет
// атрибут request_id, если он есть в контексте.
func NewContextHandler(next slog.Handler) slog.Handler {
	return contextHandler{next: next}
}

type contextHandler struct {
	next slog.Handler
}

func (h contextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if id, ok := RequestID(ctx); ok {
		r.AddAttrs(slog.String("request_id", id))
	}
	return h.next.Handle(ctx, r)
}

// WithAttrs и WithGroup обязаны снова оборачивать результат, иначе после
// logger.With(...) обёртка потеряется и request_id перестанет появляться.
func (h contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return contextHandler{next: h.next.WithAttrs(attrs)}
}

func (h contextHandler) WithGroup(name string) slog.Handler {
	return contextHandler{next: h.next.WithGroup(name)}
}
