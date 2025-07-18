package repository

import (
	"database/sql"

	"github.com/google/uuid"
	domain "github.com/liriquew/test_task/internal/domain"
)

// required because domain.User (generated code)
// does not convert IsAdmin - optional bool field
type DBUser struct {
	ID       uuid.UUID    `db:"id"`
	Username string       `db:"username"`
	Password string       `db:"password"`
	Email    string       `db:"email"`
	IsAdmin  sql.NullBool `db:"is_admin"`
}

// ConvertUserToDBUser преобразует User в DBUser.
func ConvertUserToDBUser(u domain.User) DBUser {
	dbUser := DBUser{
		ID:       uuid.UUID(u.ID.Value),
		Username: u.Username.Value,
		Password: u.Password.Value,
		Email:    u.Email.Value,
	}

	if u.IsAdmin.Set {
		dbUser.IsAdmin = sql.NullBool{
			Bool:  u.IsAdmin.Value,
			Valid: true,
		}
	}

	return dbUser
}

// ConvertDBUserToUser преобразует DBUser обратно в User.
func ConvertDBUserToUser(dbUser DBUser) domain.User {
	user := domain.User{
		ID:       domain.NewOptUUID(domain.UUID(dbUser.ID)),
		Username: domain.NewOptString(dbUser.Username),
		Password: domain.NewOptString(dbUser.Password),
		Email:    domain.NewOptString(dbUser.Email),
	}

	if dbUser.IsAdmin.Valid {
		user.IsAdmin = domain.NewOptBool(dbUser.IsAdmin.Bool)
	}

	return user
}
