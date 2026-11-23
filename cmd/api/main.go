package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/ryinner/go-auction/internal/config"
	"github.com/ryinner/go-auction/internal/handler"
	"github.com/ryinner/go-auction/internal/logger"
	"github.com/ryinner/go-auction/internal/store"
)

func main() {
	log := logger.New(os.Stdout, slog.LevelInfo)
	// для сторонних библиотек, которые пишут через slog.Default() или стандартный log
	slog.SetDefault(log)

	if err := run(log); err != nil {
		log.Error("app stopped", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	s, err := store.New(context.Background(), cfg.DBURL)
	if err != nil {
		return err
	}
	defer s.Close()

	router := handler.NewRouter(handler.Deps{
		Store:     s,
		JWTSecret: cfg.JWTSecret,
		Logger:    log,
	})

	addr := ":" + strconv.Itoa(cfg.Port)
	log.Info("http server starting", "addr", addr)

	return http.ListenAndServe(addr, router)
}
