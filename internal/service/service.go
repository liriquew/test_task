package service

import (
	"log/slog"

	"github.com/liriquew/test_task/internal/models"
)

type Repository interface {
	ListUsers() []models.User

	CreateUser(models.User) (int64, error)
	GetUserById(int64) (*models.User, error)
	UpdateUser(models.User) error
	DeleteUser(int64) error

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
