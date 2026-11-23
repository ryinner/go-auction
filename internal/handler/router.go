package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ryinner/go-auction/internal/store"
)

type Deps struct {
	Store     *store.Store
	JWTSecret string
	Logger    *slog.Logger
}

func NewRouter(d Deps) http.Handler {
	router := chi.NewRouter()

	router.Use(RequestID)
	router.Use(AccessLog(d.Logger))

	health := NewHealth(d.Store, d.Logger.With("handler", "health"))

	router.Get("/healthz", health.Check)

	return router
}
