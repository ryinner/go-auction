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

type healthResponse struct {
	Status string `json:"status"`
}

type Health struct {
	responder
	pinger pinger
}

func NewHealth(p pinger, log *slog.Logger) *Health {
	return &Health{responder: responder{log: log}, pinger: p}
}

func (h *Health) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), healthTimeout)
	defer cancel()

	if err := h.pinger.Ping(ctx); err != nil {
		h.log.ErrorContext(ctx, "db ping failed", "err", err)
		h.JSON(w, r, http.StatusServiceUnavailable, healthResponse{Status: "unavailable"})
		return
	}

	h.JSON(w, r, http.StatusOK, healthResponse{Status: "ok"})
}
