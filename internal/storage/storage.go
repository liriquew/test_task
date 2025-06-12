package storage

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/liriquew/test_task/internal/models"
)

type Storage struct {
	users     map[uuid.UUID]*models.User
	usernames map[string]uuid.UUID
	// emails    map[string]struct{}

	m *sync.RWMutex
}

func New() *Storage {
	admin := models.GetDefaultAdmin()
	usernames := map[string]uuid.UUID{
		admin.Username: admin.Id,
	}
	// emails := map[string]struct{}{
	// 	admin.Email: {},
	// }
	users := map[uuid.UUID]*models.User{
		admin.Id: admin,
	}
	return &Storage{
		users:     users,
		usernames: usernames,
		// emails:    emails,
		m: &sync.RWMutex{},
	}
}

var (
	ErrNotFound       = errors.New("user not found")
	ErrUsernameExists = errors.New("user with this username already exists")
	ErrEmailExists    = errors.New("user with this email already exists")
)

func (s *Storage) ListUsers() []models.User {
	s.m.RLock()
	defer s.m.RUnlock()

	res := make([]models.User, 0, len(s.users))

	for _, user := range s.users {
		res = append(res, *user)
	}

	return res
}

func (s *Storage) CreateUser(user models.User) (*uuid.UUID, error) {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.usernames[user.Username]; ok {
		return nil, ErrUsernameExists
	}
	// if _, ok := s.emails[user.Email]; ok {
	// 	return nil, ErrEmailExists
	// }

	var err error
	user.Id, err = uuid.NewV7()
	if err != nil {
		return nil, err
	}

	s.users[user.Id] = &user
	s.usernames[user.Username] = user.Id
	// s.emails[user.Email] = struct{}{}

	return &user.Id, nil
}

func (s *Storage) GetUserById(id uuid.UUID) (*models.User, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return nil, ErrNotFound
	}

	return user.Copy(), nil
}

func (s *Storage) UpdateUser(user models.User) error {
	s.m.Lock()
	defer s.m.Unlock()

	// is exists isn't required, user gets from Storage.GetUserById
	old := s.users[user.Id]

	if old.Username != user.Username {
		if _, ok := s.usernames[user.Username]; ok {
			return ErrUsernameExists
		}
	}

	// if old.Email != user.Email {
	// 	if _, ok := s.emails[user.Email]; ok {
	// 		return ErrEmailExists
	// 	}
	// }

	delete(s.usernames, old.Username)
	s.usernames[user.Username] = user.Id

	// delete(s.emails, old.Email)
	// s.emails[user.Email] = struct{}{}

	s.users[user.Id] = &user

	return nil
}

func (s *Storage) DeleteUser(id uuid.UUID) error {
	s.m.Lock()
	defer s.m.Unlock()
	user, ok := s.users[id]
	if !ok {
		return ErrNotFound
	}

	delete(s.users, id)
	delete(s.usernames, user.Username)
	// delete(s.emails, user.Email)

	return nil
}

func (s *Storage) GetUserByUsername(username string) (*models.User, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	idx, ok := s.usernames[username]
	if !ok {
		return nil, ErrNotFound
	}

	user := s.users[idx]

	return user.Copy(), nil
}
