package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

const healthTimeout = 2 * time.Second

type pinger interface {
	Ping(ctx context.Context) error
}

type Health struct {
	pinger pinger
	log    *slog.Logger
}

func NewHealth(p pinger, log *slog.Logger) *Health {
	return &Health{pinger: p, log: log}
}

func (h *Health) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), healthTimeout)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	if err := h.pinger.Ping(ctx); err != nil {
		h.log.ErrorContext(ctx, "db ping failed", "err", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"status":"unavailable"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
