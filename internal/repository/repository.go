package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/liriquew/test_task/internal/lib/config"
	"github.com/liriquew/test_task/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func New(cfg config.StorageConfig) *Repository {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("error while try to ping db: %w", err))
	}

	return &Repository{
		db: db,
	}
}

var (
	ErrNotFound       = errors.New("user not found")
	ErrUsernameExists = errors.New("user with this username already exists")
	ErrEmailExists    = errors.New("user with this email already exists")

	ErrEmptyUpdate = errors.New("empty fields, nothing to update")
)

func (s *Repository) ListUsers(ctx context.Context) ([]models.User, error) {
	var res []models.User

	query := `
		SELECT * FROM users
	`

	if err := s.db.SelectContext(ctx, &res, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, nil
		}

		return nil, err
	}

	return res, nil
}

func (s *Repository) CreateUser(ctx context.Context, user models.User) (*uuid.UUID, error) {
	query := `
		INSERT INTO users (username, email, password, is_admin) VALUES
		($1, $2, $3, $4) RETURNING id
	`

	err := s.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.Admin,
	).Scan(&user.Id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return nil, ErrUsernameExists
			}
		}
		return nil, err
	}

	return &user.Id, nil
}

func (s *Repository) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT * FROM users
		WHERE id = $1
	`

	user := models.User{}
	err := s.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *Repository) UpdateUser(ctx context.Context, user models.User) error {
	query := `
		UPDATE users SET %s WHERE id=$%d
	`

	queryParams, args, err := s.buildUpdate(user)
	if err != nil {
		return err
	}
	args = append(args, user.Id)

	query = fmt.Sprintf(query, queryParams, len(args))

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return ErrUsernameExists
			}
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id=$1
	`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Repository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT * FROM users
		WHERE username=$1
	`

	user := models.User{}
	err := s.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *Repository) buildUpdate(user models.User) (queryParams string, args []any, err error) {
	sb := strings.Builder{}

	if user.Username != "" {
		args = append(args, user.Username)
		sb.WriteString(fmt.Sprintf("username=$%d, ", len(args)))
	}
	if user.Email != "" {
		args = append(args, user.Email)
		sb.WriteString(fmt.Sprintf("email=$%d, ", len(args)))
	}
	if user.Admin.Valid() {
		args = append(args, user.Admin.Bool)
		sb.WriteString(fmt.Sprintf("is_admin=$%d, ", len(args)))
	}
	if len(args) == 0 {
		return "", nil, ErrEmptyUpdate
	}

	queryParams = sb.String()
	// remove last ", "
	queryParams = queryParams[:len(queryParams)-2]

	return
}
