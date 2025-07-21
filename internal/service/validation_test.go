package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		res      bool
	}{
		{
			name:     "Valid username",
			username: "username123",
			res:      true,
		},
		{
			name:     "Short username",
			username: "short",
			res:      false,
		},
		{
			name:     "Invalid character",
			username: "username;",
			res:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := validateUsername(tt.username)
			require.Equal(t, tt.res, res)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		res      bool
	}{
		{
			name:     "Valid password",
			password: "Password123",
			res:      true,
		},
		{
			name:     "Short password",
			password: "short",
			res:      false,
		},
		{
			name:     "Without lowercase",
			password: "AAAA11111",
			res:      false,
		},
		{
			name:     "Without uppercase",
			password: "aaaa11111",
			res:      false,
		},
		{
			name:     "Without digits",
			password: "aaaaAAAAA",
			res:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := validatePassword(tt.password)
			require.Equal(t, tt.res, res)
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		res   bool
	}{
		{
			name:  "Valid email",
			email: "valid@mail.ru",
			res:   true,
		},
		{
			name:  "Doesnt starts with letter",
			email: "1notvalid@mail.ru",
			res:   false,
		},
		{
			name:  "Without characters after '@'",
			email: "notvalid@.ru;",
			res:   false,
		},
		{
			name:  "Bad after '@'",
			email: "notvalid@ru",
			res:   false,
		},
		{
			name:  "Without domain",
			email: "notvalid@mail",
			res:   false,
		},
		{
			name:  "Without top level domain",
			email: "notvalid@mail.",
			res:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := validateEmail(tt.email)
			require.Equal(t, tt.res, res)
		})
	}
}
