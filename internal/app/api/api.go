package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/liriquew/test_task/internal/lib/config"
)

type Server struct {
	r chi.Router
}

type Service interface {
	ListUsers(http.ResponseWriter, *http.Request)

	CreateUser(http.ResponseWriter, *http.Request)
	GetUser(http.ResponseWriter, *http.Request)
	DeleteUser(http.ResponseWriter, *http.Request)
	PatchUser(http.ResponseWriter, *http.Request)
	PutUser(http.ResponseWriter, *http.Request)

	AuthReuqired(http.Handler) http.Handler
	CheckAdminPermission(http.Handler) http.Handler
}

func New(cfg config.AppConfig, service Service) *http.Server {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Get("/ping", Ping)
	r.With(service.AuthReuqired).Route("/users", func(r chi.Router) {
		r.Get("/", service.ListUsers)
		r.With(service.CheckAdminPermission).Route("/", func(r chi.Router) {
			r.Post("/", service.CreateUser)

			r.Route("/{userId}", func(r chi.Router) {
				r.Get("/", service.GetUser)
				r.Patch("/", service.PatchUser)
				r.Put("/", service.PutUser)
				r.Delete("/", service.DeleteUser)
			})
		})
	})

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.API.Host, cfg.API.Port),
		Handler: r,
	}
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
