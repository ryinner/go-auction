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
	auth := NewAuth(d.Store, NewTokenIssuer(d.JWTSecret), d.Logger.With("handler", "auth"))

	router.Get("/healthz", health.Check)
	router.Post("/auth/register", auth.Register)
	router.Post("/auth/login", auth.Login)

	return router
}
