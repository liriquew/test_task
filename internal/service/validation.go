package service

import (
	"regexp"
	"strings"

	domain "github.com/liriquew/test_task/internal/domain"
)

const (
	badUsername = domain.ValidationErrorMessageInvalidUsername
	badPassword = domain.ValidationErrorMessageInvalidPassword
	badEmail    = domain.ValidationErrorMessageInvalidEmail
)

func ValidateUser(user *domain.User) *domain.ValidationErrorResponse {
	if user.Username.IsSet() && !validateUsername(user.Username.Value) {
		return &domain.ValidationErrorResponse{
			Message: badUsername,
		}
	}

	if user.Password.IsSet() && !validatePassword(user.Password.Value) {
		return &domain.ValidationErrorResponse{
			Message: badPassword,
		}
	}

	if user.Email.IsSet() && !validateEmail(user.Email.Value) {
		return &domain.ValidationErrorResponse{
			Message: badEmail,
		}
	}

	return nil
}

var (
	usernameRegexp = regexp.MustCompile(`^[a-zA-Z0-9]{4,}$`)
	passwordRegexp = regexp.MustCompile(`^([a-z]|[A-Z]|[0-9]){5,}$`)
	emailRegexp    = regexp.MustCompile(`^[a-z]+[a-z0-9]*@[a-z]+\.[a-z]{2,5}$`)
)

func validateUsername(username string) bool {
	if len(username) <= 8 {
		return false
	}
	for _, char := range username {
		l := 'a' <= char && char <= 'z'
		u := 'A' <= char && char <= 'Z'
		d := '0' <= char && char <= '9'

		if !l && !u && !d {
			return false
		}
	}

	return true
}

func validatePassword(password string) bool {
	if len(password) <= 8 {
		return false
	}
	var lower, upper, digit bool

	for _, char := range password {
		l := 'a' <= char && char <= 'z'
		u := 'A' <= char && char <= 'Z'
		d := '0' <= char && char <= '9'

		if !l && !u && !d {
			return false
		}

		lower = lower || l
		upper = upper || u
		digit = digit || d
	}

	return lower && upper && digit
}

func validateEmail(email string) bool {
	email = strings.ToLower(email)
	return emailRegexp.MatchString(email)
}
