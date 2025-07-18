package app

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/liriquew/test_task/internal/lib/config"
	"github.com/liriquew/test_task/internal/service"
)

func Mux(cfg config.AppConfig, service service.Service, middmiddleware service.UserServiceMiddleware) *http.Server {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.API.Host, cfg.API.Port),
		Handler: r,
	}
}
