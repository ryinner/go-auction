package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/ryinner/go-auction/internal/config"
	"github.com/ryinner/go-auction/internal/handler"
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

	// l := slog.NewLogLogger()

	router := handler.NewRouter(handler.Deps{
		Store:     s,
		JWTSecret: cfg.JWTSecret,
	})

	return http.ListenAndServe(":"+strconv.Itoa(cfg.Port), router)
}
