package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ogen-go/ogen/middleware"

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

	// m.log.Debug("logging...", slog.String(
	// 	"vals", fmt.Sprintf("opName: %#v, t: %#v", operationName, t),
	// ))
	user, err := m.repo.GetUserByUsername(ctx, t.Username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ctx, ErrUnauthorized
		}

		m.log.Warn("error in Basic Auth", sl.Err(err))
		return nil, err
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
		m.log.Info("CheckAdminPermission", slog.Bool("isadm", isAdmin), slog.String("operationName", req.OperationName))
		if isAdmin || req.OperationName == "ServiceListUsers" {
			return next(req)
		}

		return middleware.Response{
			Type: &domain.Forbidden{},
		}, nil
	}
}

func Logging(logger *slog.Logger) middleware.Middleware {
	logger.Info("logger msg")
	return func(
		req middleware.Request,
		next func(req middleware.Request) (middleware.Response, error),
	) (middleware.Response, error) {

		return next(req)
		logger.Info("info", slog.String("operation", req.OperationName))
		logger := logger.With(
			slog.String("operation", req.OperationName),
			slog.String("operationId", req.OperationID),
		)
		logger.Info("Handling request")
		resp, err := next(req)
		if err != nil {
			logger.Error("Fail", sl.Err(err))
		} else {
			if tresp, ok := resp.Type.(interface{ GetStatusCode() int }); ok {
				logger.Info("Success",
					slog.Int("status_code", tresp.GetStatusCode()),
				)
			}
		}
		return resp, err
	}
}
