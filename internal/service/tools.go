package service

import (
	"encoding/base64"
	"fmt"

	domain "github.com/liriquew/test_task/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, *domain.InternalErrorResponse) {
	if password == "" {
		return "", nil
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", &domain.InternalErrorResponse{
			Message: domain.InternalErrorResponseMessage(
				fmt.Sprintf("error while generating password hash error: %s", err),
			),
		}
	}
	return base64.StdEncoding.EncodeToString(passwordHash), nil
}
