package handler

import (
	"context"
	"log"
	"net/http"
	"time"
)

type pinger interface {
	Ping(ctx context.Context) error
}

type Health struct {
	pinger pinger
}

func NewHealth(s pinger) *Health {
	return &Health{
		pinger: s,
	}
}

func (h *Health) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	w.Header().Set("Content-Type", "application/json")
	if err := h.pinger.Ping(ctx); err == nil {
		w.Write([]byte(`{"status": "ok"}`))
	} else {
		log.Printf("%s", err.Error())
		w.WriteHeader(503)
		w.Write([]byte(`{"status": "unavailable"}`))
	}
}
