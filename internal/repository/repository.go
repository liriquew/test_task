package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	domain "github.com/liriquew/test_task/internal/domain"
	"github.com/liriquew/test_task/internal/lib/config"

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

func (r *Repository) Close() error {
	return r.Close()
}

var (
	ErrNotFound       = errors.New("user not found")
	ErrUsernameExists = errors.New("user with this username already exists")
	ErrEmailExists    = errors.New("user with this email already exists")

	ErrEmptyUpdate = errors.New("empty fields, nothing to update")
)

const (
	usernameConstraint = "users_username_key"
	emailConstraint    = "users_email_key"
)

func UUID(id domain.UUID) string {
	return uuid.UUID(id).String()
}

func (s *Repository) ListUsers(ctx context.Context) ([]domain.User, error) {
	var users []DBUser

	query := `
		SELECT * FROM users
	`

	if err := s.db.SelectContext(ctx, &users, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.User{}, nil
		}

		return nil, err
	}

	res := make([]domain.User, 0, len(users))
	for _, user := range users {
		res = append(res, ConvertDBUserToUser(user))
	}

	return res, nil
}

func (s *Repository) CreateUser(ctx context.Context, user *domain.User) (*domain.UUID, error) {
	query := `
		INSERT INTO users (username, email, password, is_admin) VALUES
		($1, $2, $3, $4) RETURNING id
	`

	var id uuid.UUID
	err := s.db.QueryRowContext(ctx, query,
		user.Username.Value,
		user.Email.Value,
		user.Password.Value,
		user.IsAdmin.Value,
	).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				switch pqErr.Constraint {
				case usernameConstraint:
					return nil, ErrUsernameExists
				case emailConstraint:
					return nil, ErrEmailExists
				}
			}
		}
		return nil, err
	}

	res := domain.UUID(id)

	return &res, nil
}

func (s *Repository) GetUserById(ctx context.Context, id domain.UUID) (*domain.User, error) {
	query := `
		SELECT * FROM users
		WHERE id = $1
	`

	user := DBUser{}
	err := s.db.GetContext(ctx, &user, query, UUID(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	res := ConvertDBUserToUser(user)

	return &res, nil
}

func (s *Repository) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET %s WHERE id=$%d
	`

	queryParams, args, err := s.buildUpdate(user)
	if err != nil {
		return err
	}
	args = append(args, UUID(user.ID.Value))

	query = fmt.Sprintf(query, queryParams, len(args))

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				switch pqErr.Constraint {
				case usernameConstraint:
					return ErrUsernameExists
				case emailConstraint:
					return ErrEmailExists
				}
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

func (s *Repository) DeleteUser(ctx context.Context, id domain.UUID) error {
	query := `
		DELETE FROM users
		WHERE id=$1
	`

	_, err := s.db.ExecContext(ctx, query, UUID(id))
	if err != nil {
		return err
	}

	return nil
}

func (s *Repository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT * FROM users
		WHERE username=$1
	`

	user := DBUser{}
	err := s.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	res := ConvertDBUserToUser(user)

	return &res, nil
}

func (s *Repository) buildUpdate(user *domain.User) (queryParams string, args []any, err error) {
	sb := strings.Builder{}

	if user.Username.IsSet() {
		args = append(args, user.Username.Value)
		sb.WriteString(fmt.Sprintf("username=$%d, ", len(args)))
	}
	if user.Email.IsSet() {
		args = append(args, user.Email.Value)
		sb.WriteString(fmt.Sprintf("email=$%d, ", len(args)))
	}
	if user.IsAdmin.IsSet() {
		args = append(args, user.IsAdmin.Value)
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
