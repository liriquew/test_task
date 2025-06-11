package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/liriquew/test_task/internal/app/api"
	"github.com/liriquew/test_task/internal/lib/config"
	"github.com/liriquew/test_task/internal/service"
	"github.com/liriquew/test_task/internal/storage"
)

type App struct {
	srv *http.Server
	cfg config.AppConfig
}

func New(log *slog.Logger, cfg config.AppConfig) App {
	storage := storage.New()

	service := service.New(log, storage)

	server := api.New(cfg, service)

	return App{
		srv: server,
	}
}

func (s *App) Run() {
	if err := s.srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (s *App) Close(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
