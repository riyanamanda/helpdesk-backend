package user_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperror"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"
	testingutil "github.com/riyanamanda/helpdesk-backend/internal/shared/testing"
	user "github.com/riyanamanda/helpdesk-backend/internal/user"
	usermocks "github.com/riyanamanda/helpdesk-backend/internal/user/mocks"
)

func newUserService(repo *usermocks.UserRepository) user.UserService {

	return user.NewUserService(repo, config.Storage{

		PublicURL: "http://localhost:9000",

		Bucket: "helpdesk-dev",
	})

}

func TestService_CreateUser(t *testing.T) {

	testCases := []struct {
		name string

		req *user.UserCreateRequest

		setupMock func(*usermocks.UserRepository)

		assertFn func(*testing.T, error)
	}{

		{

			name: "success",

			req: &user.UserCreateRequest{

				Name: "Admin",

				Email: "  ADMIN@EMAIL.COM  ",

				Password: "password123",

				Role: user.ADMIN,

				DivisionID: 1,
			},

			setupMock: func(repo *usermocks.UserRepository) {

				repo.On("Create", mock.Anything, mock.MatchedBy(func(data *user.User) bool {

					if data == nil {

						return false

					}

					if data.Name != "Admin" {

						return false

					}

					if data.Email != "admin@email.com" {

						return false

					}

					if data.Role != user.ADMIN || data.DivisionID != 1 {

						return false

					}

					if data.Password == "password123" {

						return false

					}

					return bcrypt.CompareHashAndPassword([]byte(data.Password), []byte("password123")) == nil

				})).Run(func(args mock.Arguments) {

					data := args.Get(1).(*user.User)

					data.ID = uuid.New()

					data.CreatedAt = time.Now().UTC()

					data.UpdatedAt = data.CreatedAt

					data.IsActive = true

				}).Return(nil).Once()

			},

			assertFn: func(t *testing.T, err error) {

				require.NoError(t, err)

			},
		},

		{

			name: "already exists",

			req: &user.UserCreateRequest{

				Name: "Admin",

				Email: "admin@email.com",

				Password: "password123",

				Role: user.ADMIN,

				DivisionID: 1,
			},

			setupMock: func(repo *usermocks.UserRepository) {

				repo.On("Create", mock.Anything, mock.Anything).Return(user.ErrUserAlreadyExists).Once()

			},

			assertFn: func(t *testing.T, err error) {

				require.Error(t, err)

				testingutil.AssertAppError(t, err, apperror.CodeAlreadyExists, http.StatusConflict, "user already exists")

			},
		},

		{

			name: "repository error",

			req: &user.UserCreateRequest{

				Name: "Admin",

				Email: "admin@email.com",

				Password: "password123",

				Role: user.ADMIN,

				DivisionID: 1,
			},

			setupMock: func(repo *usermocks.UserRepository) {

				repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error")).Once()

			},

			assertFn: func(t *testing.T, err error) {

				require.Error(t, err)

				assert.EqualError(t, err, "database error")

			},
		},

		{

			name: "bcrypt error for too-long password",

			req: &user.UserCreateRequest{

				Name: "Admin",

				Email: "admin@email.com",

				Password: strings.Repeat("x", 73),

				Role: user.ADMIN,

				DivisionID: 1,
			},

			setupMock: func(_ *usermocks.UserRepository) {},

			assertFn: func(t *testing.T, err error) {

				require.Error(t, err)

				assert.Contains(t, err.Error(), "password")

			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			repo := usermocks.NewUserRepository(t)

			svc := newUserService(repo)

			tc.setupMock(repo)

			err := svc.CreateUser(context.Background(), tc.req)

			tc.assertFn(t, err)

		})

	}

}

func TestService_ListUsers(t *testing.T) {

	testCases := []struct {
		name string

		params *user.GetUserParams

		setupMock func(*usermocks.UserRepository)

		assertFn func(*testing.T, []user.UserResponse, int64, error)
	}{

		{

			name: "success with default pagination",

			params: nil,

			setupMock: func(repo *usermocks.UserRepository) {

				items := []user.UserProjection{

					{ID: uuid.New(), Name: "Admin", Email: "admin@email.com", Role: user.ADMIN, DivisionID: 1, DivisionName: "IT", IsActive: true},

					{ID: uuid.New(), Name: "Staff", Email: "staff@email.com", Role: user.EMPLOYEE, DivisionID: 2, DivisionName: "HR", IsActive: true},
				}

				repo.On("GetAll", mock.Anything, user.GetUserParams{Params: pagination.Params{Page: 1, Limit: 10}}).Return(items, int64(2), nil).Once()

			},

			assertFn: func(t *testing.T, result []user.UserResponse, total int64, err error) {

				require.NoError(t, err)

				assert.Len(t, result, 2)

				assert.Equal(t, "Admin", result[0].Name)

				assert.Equal(t, int64(2), total)

			},
		},

		{

			name: "repository error",

			params: &user.GetUserParams{},

			setupMock: func(repo *usermocks.UserRepository) {

				repo.On("GetAll", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("database error")).Once()

			},

			assertFn: func(t *testing.T, result []user.UserResponse, total int64, err error) {

				require.Error(t, err)

				assert.Empty(t, result)

				assert.Equal(t, int64(0), total)

				assert.EqualError(t, err, "database error")

			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			repo := usermocks.NewUserRepository(t)

			svc := newUserService(repo)

			tc.setupMock(repo)

			result, total, err := svc.ListUsers(context.Background(), tc.params)

			tc.assertFn(t, result, total, err)

		})

	}

}

func TestService_GetUser(t *testing.T) {

	now := time.Now().UTC()

	id := uuid.New()

	testCases := []struct {
		name string

		id *uuid.UUID

		setupMock func(*usermocks.UserRepository)

		assertFn func(*testing.T, user.UserResponse, error)
	}{

		{

			name: "success",

			id: &id,

			setupMock: func(repo *usermocks.UserRepository) {

				item := &user.UserProjection{ID: id, Name: "Admin", Email: "admin@email.com", Role: user.ADMIN, DivisionID: 1, DivisionName: "IT", IsActive: true, CreatedAt: now, UpdatedAt: now}

				repo.On("GetByID", mock.Anything, id).Return(item, nil).Once()

			},

			assertFn: func(t *testing.T, result user.UserResponse, err error) {

				require.NoError(t, err)

				assert.Equal(t, id, result.ID)

				assert.Equal(t, "Admin", result.Name)

			},
		},

		{

			name: "not found",

			id: &id,

			setupMock: func(repo *usermocks.UserRepository) {

				repo.On("GetByID", mock.Anything, id).Return(nil, user.ErrUserNotFound).Once()

			},

			assertFn: func(t *testing.T, result user.UserResponse, err error) {

				require.Error(t, err)

				assert.Equal(t, user.UserResponse{}, result)

				testingutil.AssertAppError(t, err, apperror.CodeNotFound, http.StatusNotFound, "user not found")

			},
		},

		{

			name: "repository error",

			id: &id,

			setupMock: func(repo *usermocks.UserRepository) {

				repo.On("GetByID", mock.Anything, id).Return(nil, errors.New("database error")).Once()

			},

			assertFn: func(t *testing.T, result user.UserResponse, err error) {

				require.Error(t, err)

				assert.Equal(t, user.UserResponse{}, result)

				assert.EqualError(t, err, "database error")

			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			repo := usermocks.NewUserRepository(t)

			svc := newUserService(repo)

			tc.setupMock(repo)

			result, err := svc.GetUser(context.Background(), tc.id)

			tc.assertFn(t, result, err)

		})

	}

}
