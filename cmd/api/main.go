package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ryinner/go-auction/internal/config"
	"github.com/ryinner/go-auction/internal/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()

	if err != nil {
		return err
	}

	s, err := store.New(context.Background(), cfg.DBURL)

	if err != nil {
		return err
	}

	defer s.Close()

	router := chi.NewRouter()
	router.Get("/api/index", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Hello world"}`))
	})

	return http.ListenAndServe(":"+strconv.Itoa(cfg.Port), router)
}
