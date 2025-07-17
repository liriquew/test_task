package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/liriquew/test_task/pkg/logger/sl"
	"golang.org/x/crypto/bcrypt"
)

type UserIdKey struct{}
type UserIdParam struct{}

func (s *Service) AuthReuqired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scheme := r.Header.Get("Authorization")
		if scheme == "" {
			w.Header().Add("WWW-Authenticate", `Basic realm="user service"`)
			http.Error(w, "missed authorization header", http.StatusUnauthorized)
			return
		}

		username, password, err := GetCleanCredentials(scheme)
		if err != nil {
			s.log.Warn("error while getting clean creds", sl.Err(err))
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := s.repo.GetUserByUsername(r.Context(), username)
		if err != nil {
			s.log.Warn("error while getting user by username", sl.Err(err))
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}

		hash, err := base64.StdEncoding.DecodeString(user.Password)
		if err != nil {
			s.log.Warn("error while decoding hash", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword(hash, []byte(password))
		if err != nil {
			http.Error(w, "wrong username or password", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIdKey{}, user.Id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Service) CheckAdminPermission(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetRequestUserId(r)
		user, err := s.repo.GetUserById(r.Context(), id)
		if err != nil {
			http.Error(
				w,
				"user not found, probably deleted while while request processing",
				http.StatusForbidden,
			)
			return
		}

		if !user.Admin.Bool {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetRequestUserId(r *http.Request) uuid.UUID {
	id := r.Context().Value(UserIdKey{}).(uuid.UUID)
	return id
}

func GetUserIdParam(r *http.Request) (uuid.UUID, error) {
	id := chi.URLParam(r, "userId")
	return uuid.Parse(id)
}

func GetCleanCredentials(scheme string) (username, password string, err error) {
	scheme, found := strings.CutPrefix(scheme, "Basic ")
	if !found {
		err = errors.New("wrong authorization scheme")
		return
	}

	creds, err := base64.StdEncoding.DecodeString(scheme)
	if err != nil {
		err = errors.New("bad base64 encoding")
		return
	}

	splitIdx := 0
	for ; splitIdx < len(creds); splitIdx++ {
		if creds[splitIdx] == ':' {
			break
		}
	}
	if splitIdx == len(creds) {
		err = fmt.Errorf("missed delimiter in credentials ':'")
		return
	}

	username = string(creds[:splitIdx])
	password = string(creds[splitIdx+1:])
	return
}
