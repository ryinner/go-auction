package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ryinner/go-auction/internal/store"
)

type Deps struct {
	Store     *store.Store
	JWTSecret string
}

func NewRouter(d Deps) http.Handler {
	router := chi.NewRouter()

	router.Get("/api/index", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Hello world"}`))
	})

	health := NewHealth(d.Store)

	router.Get("/healthz", health.Check)

	return router
}
