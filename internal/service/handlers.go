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

	jsontools.Encode(w, users)
}

func (s *Service) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	if err := jsontools.Decode(r.Body, &user); err != nil {
		s.log.Warn("error while decoding user in CreateUser", sl.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.repo.CreateUser(user)

	w.WriteHeader(http.StatusCreated)
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
			w.WriteHeader(http.StatusNotFound)
			return
		}

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

	s.repo.UpdateUser(*oldUser)

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

	s.repo.UpdateUser(newUser)

	w.WriteHeader(http.StatusOK)
}
