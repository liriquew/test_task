package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"errors"

	"github.com/go-chi/chi/v5"
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
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := s.repo.GetUserByUsername(username)
		if err != nil {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}

		if user.Password != password {
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
		user, err := s.repo.GetUserById(id)
		if err != nil {
			http.Error(
				w,
				"user not found, probably deleted while while request processing",
				http.StatusForbidden,
			)
			return
		}

		if !user.Admin.Value() {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetRequestUserId(r *http.Request) int64 {
	id := r.Context().Value(UserIdKey{}).(int64)
	return id
}

func GetUserIdParam(r *http.Request) (id int64, err error) {
	id, err = strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	return
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
