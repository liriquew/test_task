package service

import (
	"context"
	"log/slog"

	domain "github.com/liriquew/test_task/internal/domain"
)

//go:generate mockgen -source=service.go -destination=mocks/repository.go -package=mocks
type Repository interface {
	ListUsers(context.Context) ([]domain.User, error)

	CreateUser(context.Context, *domain.User) (*domain.UUID, error)
	GetUserById(context.Context, domain.UUID) (*domain.User, error)
	UpdateUser(context.Context, *domain.User) error
	DeleteUser(context.Context, domain.UUID) error

	GetUserByUsername(context.Context, string) (*domain.User, error)
}

type Service struct {
	repo Repository
	log  *slog.Logger
}

func New(log *slog.Logger, repo Repository) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

type UserServiceMiddleware struct {
	log  *slog.Logger
	repo Repository
}

func NewMiddleware(log *slog.Logger, repo Repository) *UserServiceMiddleware {
	return &UserServiceMiddleware{
		log:  log,
		repo: repo,
	}
}
