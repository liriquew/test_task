package service_test

import (
	"context"
	"log/slog"
	"testing"

	domain "github.com/liriquew/test_task/internal/domain"
	"github.com/liriquew/test_task/internal/repository"
	"github.com/liriquew/test_task/internal/service"
	"github.com/liriquew/test_task/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func StubLogger() *slog.Logger {
	log := slog.New(slog.DiscardHandler)
	return log
}

func TestListUsers(t *testing.T) {
	t.Parallel()
	repo := mocks.NewMockRepository(gomock.NewController(t))

	users := []domain.User{
		{
			ID:       domain.NewOptUUID(domain.UUID{}),
			Username: domain.NewOptString("username"),
			Password: domain.NewOptString("password"),
			Email:    domain.NewOptString("email"),
		},
	}
	res := domain.ServiceListUsersOKApplicationJSON(users)
	repo.
		EXPECT().
		ListUsers(gomock.Any()).
		Return(users, nil)
	s := service.New(StubLogger(), repo)

	resp, err := s.ServiceListUsers(context.Background())

	require.Nil(t, err)
	usersResp, ok := resp.(*domain.ServiceListUsersOKApplicationJSON)
	require.True(t, ok)
	require.Equal(t,
		&res,
		usersResp,
	)
}

func TestCreateUser(t *testing.T) {
	type deps struct {
		repo *mocks.MockRepository
	}

	type test struct {
		name    string
		setup   func(d deps, t *test)
		user    domain.User
		res     domain.ServiceCreateUserRes
		wantErr bool
	}

	tests := []test{
		{
			name: "All valid",
			setup: func(d deps, t *test) {
				d.repo.EXPECT().
					CreateUser(gomock.Any(), &t.user).
					Return(&domain.UUID{}, nil)
			},
			user: domain.User{
				Username: domain.NewOptString("username1"),
				Password: domain.NewOptString("password123A"),
				Email:    domain.NewOptString("valid@mail.ru"),
			},
			res: &domain.User{
				ID:       domain.NewOptUUID(domain.UUID{}),
				Username: domain.NewOptString("username1"),
				Password: domain.NewOptString(""),
				Email:    domain.NewOptString("valid@mail.ru"),
			},
			wantErr: false,
		},
		{
			name:  "Invalid Username",
			setup: nil,
			user: domain.User{
				Username: domain.NewOptString("юзернейм1"),
				Password: domain.NewOptString("password123A"),
				Email:    domain.NewOptString("valid@mail.ru"),
			},
			res: &domain.ValidationErrorResponse{
				Message: domain.ValidationErrorMessageInvalidUsername,
			},
			wantErr: false,
		},
		{
			name:  "Invalid Password",
			setup: nil,
			user: domain.User{
				Username: domain.NewOptString("username1"),
				// without upper character
				Password: domain.NewOptString("password123password"),
				Email:    domain.NewOptString("valid@mail.ru"),
			},
			res: &domain.ValidationErrorResponse{
				Message: domain.ValidationErrorMessageInvalidPassword,
			},
			wantErr: false,
		},
		{
			name:  "Invalid Email",
			setup: nil,
			user: domain.User{
				Username: domain.NewOptString("username1"),
				Password: domain.NewOptString("password123A"),
				Email:    domain.NewOptString("@mail.ru"),
			},
			res: &domain.ValidationErrorResponse{
				Message: domain.ValidationErrorMessageInvalidEmail,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d := deps{
				repo: mocks.NewMockRepository(ctrl),
			}
			if tt.setup != nil {
				tt.setup(d, &tt)
			}
			s := service.New(StubLogger(), d.repo)

			res, err := s.ServiceCreateUser(context.Background(), &tt.user)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, res)
				return
			}

			if v, ok := res.(*domain.User); ok {
				v.Password.Value = ""
			}
			require.Equal(t, res, tt.res)
		})
	}
}

func TestPatchUser(t *testing.T) {
	type deps struct {
		repo *mocks.MockRepository
	}

	type test struct {
		name    string
		setup   func(d deps, t *test)
		user    domain.User
		res     domain.ServicePatchUserRes
		wantErr bool
	}

	tests := []test{
		{
			name: "All valid",
			setup: func(d deps, t *test) {
				d.repo.EXPECT().
					UpdateUser(gomock.Any(), &t.user).
					Return(nil)
			},
			user: domain.User{
				Username: domain.NewOptString("username1"),
				Password: domain.NewOptString("password123A"),
				Email:    domain.NewOptString("valid@mail.ru"),
			},
			res:     &domain.ServicePatchUserOK{},
			wantErr: false,
		},
		{
			name:  "Invalid Username",
			setup: nil,
			user: domain.User{
				Username: domain.NewOptString("юзернейм1"),
				Password: domain.NewOptString("password123A"),
				Email:    domain.NewOptString("valid@mail.ru"),
			},
			res: &domain.ValidationErrorResponse{
				Message: domain.ValidationErrorMessageInvalidUsername,
			},
			wantErr: false,
		},
		{
			name:  "Invalid Password",
			setup: nil,
			user: domain.User{
				Username: domain.NewOptString("username1"),
				// without upper character
				Password: domain.NewOptString("password123password"),
				Email:    domain.NewOptString("valid@mail.ru"),
			},
			res: &domain.ValidationErrorResponse{
				Message: domain.ValidationErrorMessageInvalidPassword,
			},
			wantErr: false,
		},
		{
			name:  "Invalid Email",
			setup: nil,
			user: domain.User{
				Username: domain.NewOptString("username1"),
				Password: domain.NewOptString("password123A"),
				Email:    domain.NewOptString("@mail.ru"),
			},
			res: &domain.ValidationErrorResponse{
				Message: domain.ValidationErrorMessageInvalidEmail,
			},
			wantErr: false,
		},
		{
			name: "Empty Update",
			setup: func(d deps, t *test) {
				d.repo.EXPECT().
					UpdateUser(gomock.Any(), &t.user).
					Return(repository.ErrEmptyUpdate)
			},
			user: domain.User{},
			res: &domain.ValidationErrorResponse{
				Message: "nothing to update",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d := deps{
				repo: mocks.NewMockRepository(ctrl),
			}
			if tt.setup != nil {
				tt.setup(d, &tt)
			}
			s := service.New(StubLogger(), d.repo)

			res, err := s.ServicePatchUser(context.Background(), &tt.user, domain.ServicePatchUserParams{
				UserId: domain.UUID{},
			})
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, res)
				return
			}

			require.Equal(t, res, tt.res)
		})
	}
}

func TestPutUser(t *testing.T) {
	type deps struct {
		repo *mocks.MockRepository
	}

	type test struct {
		name    string
		setup   func(d deps, t *test)
		user    domain.User
		res     domain.ServicePutUserRes
		wantErr bool
	}

	tests := []test{
		{
			name: "All valid",
			setup: func(d deps, t *test) {
				d.repo.EXPECT().
					UpdateUser(gomock.Any(), &t.user).
					Return(nil)
			},
			user: domain.User{
				Username: domain.NewOptString("username1"),
				Password: domain.NewOptString("password123A"),
				Email:    domain.NewOptString("valid@mail.ru"),
			},
			res:     &domain.ServicePutUserOK{},
			wantErr: false,
		},
		{
			name:  "Invalid Username",
			setup: nil,
			user: domain.User{
				Username: domain.NewOptString("юзернейм1"),
				Password: domain.NewOptString("password123A"),
				Email:    domain.NewOptString("valid@mail.ru"),
			},
			res: &domain.ValidationErrorResponse{
				Message: domain.ValidationErrorMessageInvalidUsername,
			},
			wantErr: false,
		},
		{
			name:  "Invalid Password",
			setup: nil,
			user: domain.User{
				Username: domain.NewOptString("username1"),
				// without upper character
				Password: domain.NewOptString("password123password"),
				Email:    domain.NewOptString("valid@mail.ru"),
			},
			res: &domain.ValidationErrorResponse{
				Message: domain.ValidationErrorMessageInvalidPassword,
			},
			wantErr: false,
		},
		{
			name:  "Invalid Email",
			setup: nil,
			user: domain.User{
				Username: domain.NewOptString("username1"),
				Password: domain.NewOptString("password123A"),
				Email:    domain.NewOptString("@mail.ru"),
			},
			res: &domain.ValidationErrorResponse{
				Message: domain.ValidationErrorMessageInvalidEmail,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d := deps{
				repo: mocks.NewMockRepository(ctrl),
			}
			if tt.setup != nil {
				tt.setup(d, &tt)
			}
			s := service.New(StubLogger(), d.repo)

			res, err := s.ServicePutUser(context.Background(), &tt.user, domain.ServicePutUserParams{
				UserId: domain.UUID{},
			})
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, res)
				return
			}

			require.Equal(t, res, tt.res)
		})
	}
}
