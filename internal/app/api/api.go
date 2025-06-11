package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

func New(service Service) *http.Server {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.With(service.AuthReuqired).Route("/users", func(r chi.Router) {
		r.Get("/", service.ListUsers)
		r.Get("/{userId}", service.GetUser)

		r.With(service.CheckAdminPermission).Route("/", func(r chi.Router) {
			r.Post("/", service.CreateUser)

			r.Route("/{userId}", func(r chi.Router) {
				r.Patch("/", service.PatchUser)
				r.Put("/", service.PutUser)
				r.Delete("/", service.DeleteUser)
			})
		})
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}
