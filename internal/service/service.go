package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/liriquew/test_task/internal/models"
)

type Repository interface {
	ListUsers(context.Context) ([]models.User, error)

	CreateUser(context.Context, models.User) (*uuid.UUID, error)
	GetUserById(context.Context, uuid.UUID) (*models.User, error)
	UpdateUser(context.Context, models.User) error
	DeleteUser(context.Context, uuid.UUID) error

	GetUserByUsername(context.Context, string) (*models.User, error)
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
