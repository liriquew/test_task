package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	domain "github.com/liriquew/test_task/internal/domain"
	"github.com/liriquew/test_task/internal/repository"
	"github.com/liriquew/test_task/pkg/logger/sl"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Health(ctx context.Context) error {
	return nil
}

func (s *Service) ServiceListUsers(
	ctx context.Context,
	params domain.ServiceListUsersParams,
) (
	domain.ServiceListUsersRes,
	error,
) {
	users, err := s.repo.ListUsers(ctx, params.Offset.Value)
	if err != nil {
		s.log.Warn("error while getting users in ListUsers", sl.Err(err))
		return &domain.InternalErrorResponse{}, nil
	}

	for i := range users {
		users[i].Password.Value = ""
	}

	res := domain.ServiceListUsersOKApplicationJSON(users)

	return &res, nil
}

func (s *Service) ServiceCreateUser(
	ctx context.Context,
	user *domain.User,
) (domain.ServiceCreateUserRes, error) {
	if user.Username.Value == "" {
		return &domain.ValidationErrorResponse{
			Message: "empty username",
		}, nil
	}
	if user.Password.Value == "" {
		return &domain.ValidationErrorResponse{
			Message: "empty password",
		}, nil
	}
	if user.Email.Value == "" {
		return &domain.ValidationErrorResponse{
			Message: "empty email",
		}, nil
	}

	// validate user
	if errResp := ValidateUser(user); errResp != nil {
		return errResp, nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password.Value), bcrypt.DefaultCost)
	if err != nil {
		s.log.Warn("error while generating password hash", sl.Err(err))
		return &domain.InternalErrorResponse{
			Message: domain.InternalErrorResponseMessage(
				fmt.Sprintf("error while generating password hash error: %s", err),
			),
		}, nil
	}
	user.Password.Value = base64.StdEncoding.EncodeToString(passwordHash)

	uuid, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		s.log.Warn("error while creating user", sl.Err(err))
		if errors.Is(err, repository.ErrUsernameExists) {
			return &domain.AlreadyExistsResponse{
				Message: "username already exists",
			}, nil
		}
		if errors.Is(err, repository.ErrEmailExists) {
			return &domain.AlreadyExistsResponse{
				Message: "email already exists",
			}, nil
		}

		return &domain.InternalErrorResponse{
			Message: domain.InternalErrorResponseMessage(
				fmt.Sprintf("internal error while creating user error: %s", err),
			),
		}, nil
	}

	user.ID.SetTo(*uuid)

	return user, nil
}

func (s *Service) ServiceGetUser(
	ctx context.Context,
	params domain.ServiceGetUserParams,
) (domain.ServiceGetUserRes, error) {
	user, err := s.repo.GetUserById(ctx, params.UserId)
	if err != nil {
		s.log.Warn("error while getting user by id", sl.Err(err))
		if errors.Is(err, repository.ErrNotFound) {
			return &domain.NotFoundResponse{
				Message: "user not found",
			}, nil
		}

		return &domain.InternalErrorResponse{
			Message: domain.InternalErrorResponseMessage(
				fmt.Sprintf("internal error: %s", err),
			),
		}, nil
	}

	return user, nil
}

func (s *Service) ServiceDeleteUser(
	ctx context.Context,
	params domain.ServiceDeleteUserParams,
) (domain.ServiceDeleteUserRes, error) {
	err := s.repo.DeleteUser(ctx, params.UserId)
	if err != nil {
		s.log.Warn("error while deleting user", sl.Err(err))
		return &domain.InternalErrorResponse{
			Message: domain.InternalErrorResponseMessage(
				fmt.Sprintf("internal error: %s", err),
			),
		}, nil
	}

	return &domain.ServiceDeleteUserOK{}, nil
}

func (s *Service) ServicePatchUser(
	ctx context.Context,
	user *domain.User,
	params domain.ServicePatchUserParams,
) (domain.ServicePatchUserRes, error) {
	user.ID.Value = params.UserId

	// validate user
	if errResp := ValidateUser(user); errResp != nil {
		return errResp, nil
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		s.log.Warn("error while patching user PatchUser", sl.Err(err))
		if errors.Is(err, repository.ErrUsernameExists) {
			return &domain.AlreadyExistsResponse{
				Message: "username already exists",
			}, nil
		}
		if errors.Is(err, repository.ErrEmailExists) {
			return &domain.AlreadyExistsResponse{
				Message: "email already exists",
			}, nil
		}
		if errors.Is(err, repository.ErrEmptyUpdate) {
			return &domain.ValidationErrorResponse{
				Message: "nothing to update",
			}, nil
		}

		return &domain.InternalErrorResponse{
			Message: domain.InternalErrorResponseMessage(
				fmt.Sprintf("internal error: %s", err),
			),
		}, nil
	}

	return &domain.ServicePatchUserOK{}, nil
}

func (s *Service) ServicePutUser(
	ctx context.Context,
	user *domain.User,
	params domain.ServicePutUserParams,
) (domain.ServicePutUserRes, error) {
	if user.Username.Value == "" {
		return &domain.ValidationErrorResponse{
			Message: "empty username",
		}, nil
	}
	if user.Password.Value == "" {
		return &domain.ValidationErrorResponse{
			Message: "empty username",
		}, nil
	}
	if user.Email.Value == "" {
		return &domain.ValidationErrorResponse{
			Message: "empty username",
		}, nil
	}

	// validate user
	if errResp := ValidateUser(user); errResp != nil {
		return errResp, nil
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		s.log.Warn("error while updating user in PutUser", sl.Err(err))
		if errors.Is(err, repository.ErrUsernameExists) {
			return &domain.ValidationErrorResponse{
				Message: "username already exists",
			}, nil
		}
		if errors.Is(err, repository.ErrEmailExists) {
			return &domain.ValidationErrorResponse{
				Message: "email already exists",
			}, nil
		}

		return &domain.InternalErrorResponse{
			Message: domain.InternalErrorResponseMessage(
				fmt.Sprintf("internal error: %s", err),
			),
		}, nil
	}

	return &domain.ServicePutUserOK{}, nil
}
