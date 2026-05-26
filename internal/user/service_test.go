package user_test

import (
	"context"
	"errors"
	"io"
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
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"
	testingutil "github.com/riyanamanda/helpdesk-backend/internal/shared/testing"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/upload"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
	user "github.com/riyanamanda/helpdesk-backend/internal/user"
	usermocks "github.com/riyanamanda/helpdesk-backend/internal/user/mocks"
)

// MockStorage provides a no-op implementation of storage.Storage for testing
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	args := m.Called(ctx, key, reader, size, contentType)
	return args.Error(0)
}

func (m *MockStorage) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func newMockStorage() *MockStorage {
	store := new(MockStorage)
	return store
}

func TestService_RegisterUser(t *testing.T) {
	testCases := []struct {
		name      string
		req       *user.UserCreateRequest
		setupMock func(*usermocks.UserRepository)
		assertFn  func(*testing.T, error)
	}{
		{
			name: "success",
			req: &user.UserCreateRequest{
				Name:       "Admin",
				Email:      "  ADMIN@EMAIL.COM  ",
				Password:   "password123",
				Role:       user.ADMIN,
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
				Name:       "Admin",
				Email:      "admin@email.com",
				Password:   "password123",
				Role:       user.ADMIN,
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
				Name:       "Admin",
				Email:      "admin@email.com",
				Password:   "password123",
				Role:       user.ADMIN,
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
				Name:       "Admin",
				Email:      "admin@email.com",
				Password:   strings.Repeat("x", 73),
				Role:       user.ADMIN,
				DivisionID: 1,
			},
			setupMock: func(repo *usermocks.UserRepository) {
			},
			assertFn: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "password")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := usermocks.NewUserRepository(t)
			storage := newMockStorage()
			svc := user.NewUserService(repo, storage, config.Storage{
				PublicURL: "http://localhost:9000",
				Bucket:    "helpdesk-dev",
			})
			tc.setupMock(repo)

			err := svc.RegisterUser(context.Background(), tc.req)
			tc.assertFn(t, err)
		})
	}
}

func TestService_FetchAllUsers(t *testing.T) {
	testCases := []struct {
		name      string
		params    *user.GetUserParams
		setupMock func(*usermocks.UserRepository)
		assertFn  func(*testing.T, []user.UserResponse, int64, error)
	}{
		{
			name:   "success with default pagination",
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
			name:   "repository error",
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
			storage := newMockStorage()
			svc := user.NewUserService(repo, storage, config.Storage{
				PublicURL: "http://localhost:9000",
				Bucket:    "helpdesk-dev",
			})
			tc.setupMock(repo)

			result, total, err := svc.FetchAllUsers(context.Background(), tc.params)
			tc.assertFn(t, result, total, err)
		})
	}
}

func TestService_FindUserByID(t *testing.T) {
	now := time.Now().UTC()
	id := uuid.New()

	testCases := []struct {
		name      string
		id        *uuid.UUID
		setupMock func(*usermocks.UserRepository)
		assertFn  func(*testing.T, user.UserResponse, error)
	}{
		{
			name: "success",
			id:   &id,
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
			id:   &id,
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
			id:   &id,
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
			storage := newMockStorage()
			svc := user.NewUserService(repo, storage, config.Storage{
				PublicURL: "http://localhost:9000",
				Bucket:    "helpdesk-dev",
			})
			tc.setupMock(repo)

			result, err := svc.FindUserByID(context.Background(), tc.id)
			tc.assertFn(t, result, err)
		})
	}
}

func newUploadFile(content string) *upload.File {
	return &upload.File{
		Content:     strings.NewReader(content),
		Filename:    "avatar.png",
		ContentType: "image/png",
		Size:        1024,
	}
}

func TestService_UpdateUserAvatar(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := usermocks.NewUserRepository(t)
		storage := newMockStorage()
		svc := user.NewUserService(repo, storage, config.Storage{
			PublicURL: "http://localhost:9000",
			Bucket:    "helpdesk-dev",
		})

		userID := uuid.New()
		ctx := utils.SetUserIDToContext(context.Background(), userID)

		storage.On("Upload", mock.Anything, "avatars/"+userID.String()+"/avatar", mock.Anything, int64(1024), "image/png").Return(nil).Once()
		repo.On("UpdateAvatar", mock.Anything, userID, "avatars/"+userID.String()+"/avatar").Return(nil).Once()

		err := svc.UpdateUserAvatar(ctx, newUploadFile("fake image content"))
		require.NoError(t, err)
	})

	t.Run("unauthorized when user id not in context", func(t *testing.T) {
		repo := usermocks.NewUserRepository(t)
		storage := newMockStorage()
		svc := user.NewUserService(repo, storage, config.Storage{
			PublicURL: "http://localhost:9000",
			Bucket:    "helpdesk-dev",
		})

		err := svc.UpdateUserAvatar(context.Background(), newUploadFile("fake image content"))
		require.Error(t, err)
		testingutil.AssertAppError(t, err, apperror.CodeForbidden, http.StatusForbidden, "unauthorized")
	})

	t.Run("storage upload error", func(t *testing.T) {
		repo := usermocks.NewUserRepository(t)
		storage := newMockStorage()
		svc := user.NewUserService(repo, storage, config.Storage{
			PublicURL: "http://localhost:9000",
			Bucket:    "helpdesk-dev",
		})

		userID := uuid.New()
		ctx := utils.SetUserIDToContext(context.Background(), userID)

		storage.On("Upload", mock.Anything, "avatars/"+userID.String()+"/avatar", mock.Anything, int64(1024), "image/png").Return(errors.New("storage error")).Once()

		err := svc.UpdateUserAvatar(ctx, newUploadFile("fake image content"))
		require.Error(t, err)
		assert.EqualError(t, err, "storage error")
	})

	t.Run("repository update avatar error", func(t *testing.T) {
		repo := usermocks.NewUserRepository(t)
		storage := newMockStorage()
		svc := user.NewUserService(repo, storage, config.Storage{
			PublicURL: "http://localhost:9000",
			Bucket:    "helpdesk-dev",
		})

		userID := uuid.New()
		ctx := utils.SetUserIDToContext(context.Background(), userID)

		storage.On("Upload", mock.Anything, "avatars/"+userID.String()+"/avatar", mock.Anything, int64(1024), "image/png").Return(nil).Once()
		repo.On("UpdateAvatar", mock.Anything, userID, "avatars/"+userID.String()+"/avatar").Return(errors.New("repository error")).Once()

		err := svc.UpdateUserAvatar(ctx, newUploadFile("fake image content"))
		require.Error(t, err)
		assert.EqualError(t, err, "repository error")
	})
}
