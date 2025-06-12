package service

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/liriquew/test_task/internal/models"
)

type Repository interface {
	ListUsers() []models.User

	CreateUser(models.User) (*uuid.UUID, error)
	GetUserById(uuid.UUID) (*models.User, error)
	UpdateUser(models.User) error
	DeleteUser(uuid.UUID) error

	GetUserByUsername(string) (*models.User, error)
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
