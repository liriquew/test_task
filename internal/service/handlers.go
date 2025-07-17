package service

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/liriquew/test_task/internal/lib/jsontools"
	"github.com/liriquew/test_task/internal/models"
	"github.com/liriquew/test_task/internal/repository"
	"github.com/liriquew/test_task/pkg/logger/sl"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.repo.ListUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Warn("error while generating password hash", sl.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = base64.StdEncoding.EncodeToString(passwordHash)

	id, err := s.repo.CreateUser(r.Context(), user)
	if err != nil {
		s.log.Warn("error while creating user", sl.Err(err))
		if errors.Is(err, repository.ErrUsernameExists) {
			http.Error(w, "username already taken", http.StatusConflict)
			return
		}
		if errors.Is(err, repository.ErrEmailExists) {
			http.Error(w, "email already taken", http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := struct {
		Id uuid.UUID `json:"id"`
	}{
		Id: *id,
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

	user, err := s.repo.GetUserById(r.Context(), userId)
	if err != nil {
		s.log.Warn("error while getting user by id", sl.Err(err))
		if errors.Is(err, repository.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

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

	err = s.repo.DeleteUser(r.Context(), userId)
	if err != nil {
		s.log.Warn("error while deleting user", sl.Err(err))
		if errors.Is(err, repository.ErrNotFound) {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) PatchUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	if err := jsontools.Decode(r.Body, &user); err != nil {
		s.log.Warn("error while decoding body in PatchUser", sl.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := GetUserIdParam(r)
	if err != nil {
		s.log.Warn("error while getting userId param in DeleteUser", sl.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.Id = userId

	if err := s.repo.UpdateUser(r.Context(), user); err != nil {
		s.log.Warn("error while patching user PatchUser", sl.Err(err))
		if errors.Is(err, repository.ErrUsernameExists) {
			http.Error(w, "username already taken", http.StatusConflict)
			return
		}
		if errors.Is(err, repository.ErrEmailExists) {
			http.Error(w, "email already taken", http.StatusConflict)
			return
		}

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

	if err := s.repo.UpdateUser(r.Context(), newUser); err != nil {
		s.log.Warn("error while updating user in PutUser", sl.Err(err))
		if errors.Is(err, repository.ErrUsernameExists) {
			http.Error(w, "username already taken", http.StatusConflict)
			return
		}
		if errors.Is(err, repository.ErrEmailExists) {
			http.Error(w, "email already taken", http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
