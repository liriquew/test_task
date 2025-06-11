package service

import (
	"errors"
	"net/http"

	"log/slog"

	"github.com/liriquew/test_task/internal/lib/jsontools"
	"github.com/liriquew/test_task/internal/models"
	"github.com/liriquew/test_task/internal/storage"
	"github.com/liriquew/test_task/pkg/logger/sl"
)

func (s *Service) ListUsers(w http.ResponseWriter, r *http.Request) {
	users := s.repo.ListUsers()

	s.log.Info("ListUsers:", slog.Any(
		"users", users,
	))
	for i := range users {
		users[i].Password = ""
	}

	jsontools.Encode(w, users)
}

func (s *Service) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	if err := jsontools.Decode(r.Body, &user); err != nil {
		s.log.Warn("error while decoding user in CreateUser", sl.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(w, "empty username", http.StatusBadRequest)
		return
	}
	if user.Password == "" {
		http.Error(w, "empty password", http.StatusBadRequest)
		return
	}
	if user.Email == "" {
		http.Error(w, "empty email", http.StatusBadRequest)
		return
	}

	id, err := s.repo.CreateUser(user)
	if err != nil {
		if errors.Is(err, storage.ErrUsernameExists) {
			http.Error(w, "username already taken", http.StatusConflict)
			return
		}
		if errors.Is(err, storage.ErrEmailExists) {
			http.Error(w, "email already taken", http.StatusConflict)
			return
		}

		// never happen
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := struct {
		Id int64 `json:"id"`
	}{
		Id: id,
	}
	w.WriteHeader(http.StatusCreated)
	jsontools.Encode(w, resp)
}

func (s *Service) GetUser(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserIdParam(r)
	if err != nil {
		s.log.Warn("error while getting userId param in GetUser", sl.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := s.repo.GetUserById(userId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		s.log.Warn("error while getting user in GetUser", sl.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Password = ""

	jsontools.Encode(w, user)
}

func (s *Service) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserIdParam(r)
	if err != nil {
		s.log.Warn("error while getting userId param in DeleteUser", sl.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.repo.DeleteUser(userId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			w.WriteHeader(http.StatusOK)
			return
		}

		// never happen
		s.log.Warn("error while getting user in DeleteUser", sl.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) PatchUser(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserIdParam(r)
	if err != nil {
		s.log.Warn("error while getting userId param in DeleteUser", sl.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	oldUser, err := s.repo.GetUserById(userId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// never happen
		s.log.Warn("error while patch user", sl.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newUser := models.User{}
	if err := jsontools.Decode(r.Body, &newUser); err != nil {
		s.log.Warn("error while decoding body in PatchUser", sl.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	oldUser.Patch(newUser)

	if err := s.repo.UpdateUser(*oldUser); err != nil {
		if errors.Is(err, storage.ErrUsernameExists) {
			http.Error(w, "username already taken", http.StatusConflict)
			return
		}
		if errors.Is(err, storage.ErrEmailExists) {
			http.Error(w, "email already taken", http.StatusConflict)
			return
		}

		// never happen
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) PutUser(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserIdParam(r)
	if err != nil {
		s.log.Warn("error while getting userId param in DeleteUser", sl.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newUser := models.User{
		Id: userId,
	}
	if err := jsontools.Decode(r.Body, &newUser); err != nil {
		s.log.Warn("error while decoding body in PatchUser", sl.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if newUser.Username == "" {
		http.Error(w, "empty username", http.StatusBadRequest)
		return
	}
	if newUser.Password == "" {
		http.Error(w, "empty password", http.StatusBadRequest)
		return
	}
	if newUser.Email == "" {
		http.Error(w, "empty email", http.StatusBadRequest)
		return
	}

	if err := s.repo.UpdateUser(newUser); err != nil {
		if errors.Is(err, storage.ErrUsernameExists) {
			http.Error(w, "username already taken", http.StatusConflict)
			return
		}
		if errors.Is(err, storage.ErrEmailExists) {
			http.Error(w, "email already taken", http.StatusConflict)
			return
		}

		// never happen
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
