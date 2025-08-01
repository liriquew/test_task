package service

import (
	"context"
	"encoding/base64"
	"errors"
	"log/slog"

	"github.com/ogen-go/ogen/middleware"
	"golang.org/x/crypto/bcrypt"

	api "github.com/liriquew/test_task/internal/domain"
	domain "github.com/liriquew/test_task/internal/domain"
	"github.com/liriquew/test_task/internal/repository"
	"github.com/liriquew/test_task/pkg/logger/sl"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type (
	IsAdmin struct{}
)

func (m *UserServiceMiddleware) HandleBasicAuth(
	ctx context.Context,
	operationName api.OperationName,
	t domain.BasicAuth,
) (context.Context, error) {
	user, err := m.repo.GetUserByUsername(ctx, t.Username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ctx, ErrUnauthorized
		}

		m.log.Warn("error in Basic Auth", sl.Err(err))
		return nil, err
	}

	passwordHash, err := base64.StdEncoding.DecodeString(user.Password.Value)
	if err != nil {
		m.log.Warn("error while decoding password hash", sl.Err(err))
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(passwordHash, []byte(t.Password)); err != nil {
		return ctx, ErrUnauthorized
	}

	ctx = context.WithValue(ctx, IsAdmin{}, user.IsAdmin.Value)
	return ctx, nil
}

func (m *UserServiceMiddleware) CheckAdminPermission() middleware.Middleware {
	return func(
		req middleware.Request,
		next middleware.Next,
	) (middleware.Response, error) {
		isAdmin := req.Context.Value(IsAdmin{}).(bool)
		if isAdmin || req.OperationName == "ServiceListUsers" {
			return next(req)
		}

		return middleware.Response{
			Type: &domain.ForbiddenResponse{},
		}, nil
	}
}

func Logging(logger *slog.Logger) middleware.Middleware {
	logger.Info("logger msg")
	return func(
		req middleware.Request,
		next func(req middleware.Request) (middleware.Response, error),
	) (middleware.Response, error) {
		attrs := []any{
			slog.String("operation", req.OperationName),
		}
		logger := logger.With(
			slog.String("operation", req.OperationName),
			slog.String("operationId", req.OperationID),
		)
		resp, err := next(req)
		if err != nil {
			attrs = append(attrs, sl.Err(err))
		} else {
			if tresp, ok := resp.Type.(interface{ GetStatusCode() int }); ok {
				attrs = append(attrs, slog.Int("status_code", tresp.GetStatusCode()))
			}
		}

		logger.Info("query", attrs...)

		return resp, err
	}
}
