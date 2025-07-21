package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	domain "github.com/liriquew/test_task/internal/domain"
	"github.com/liriquew/test_task/internal/lib/config"
	"github.com/liriquew/test_task/internal/repository"
	"github.com/liriquew/test_task/internal/service"
)

type App struct {
	srv     *http.Server
	closers []func() error
}

func New(log *slog.Logger, cfg config.AppConfig) App {
	storage := repository.New(cfg.Storage)
	srvs := service.New(log, storage)
	mdlwr := service.NewMiddleware(log, storage)

	server, err := domain.NewServer(srvs, mdlwr, []domain.ServerOption{
		domain.WithMiddleware(
			service.Logging(log),
			mdlwr.CheckAdminPermission(),
		),
	}...)
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf("%s:%d", cfg.API.Host, cfg.API.Port)

	return App{
		srv: &http.Server{
			Handler: server,
			Addr:    addr,
		},
		closers: []func() error{
			storage.Close,
		},
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
