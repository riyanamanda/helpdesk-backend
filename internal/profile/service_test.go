package profile_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	profile "github.com/riyanamanda/helpdesk-backend/internal/profile"
	profilemocks "github.com/riyanamanda/helpdesk-backend/internal/profile/mocks"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	testingutil "github.com/riyanamanda/helpdesk-backend/internal/shared/testing"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
	storagemocks "github.com/riyanamanda/helpdesk-backend/internal/storage/mocks"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func newProfileService(repo profile.ProfileRepository, storage *storagemocks.Storage) profile.ProfileService {
	return profile.NewProfileService(repo, storage, config.Storage{}, config.Auth{})
}

func ctxWithUser(id uuid.UUID) context.Context {
	return utils.SetUserIDToContext(context.Background(), id)
}

func TestService_GetProfile(t *testing.T) {
	userID := uuid.New()

	testCases := []struct {
		name      string
		ctx       context.Context
		setupMock func(*profilemocks.ProfileRepository)
		assertFn  func(*testing.T, profile.ProfileResponse, error)
	}{
		{
			name: "success",
			ctx:  ctxWithUser(userID),
			setupMock: func(repo *profilemocks.ProfileRepository) {
				proj := &user.UserProjection{ID: userID, Name: "Alice", Email: "alice@example.com"}
				repo.On("GetByID", mock.Anything, userID).Return(proj, nil).Once()
			},
			assertFn: func(t *testing.T, result profile.ProfileResponse, err error) {
				require.NoError(t, err)
				assert.Equal(t, "Alice", result.Name)
			},
		},
		{
			name:      "no user in context",
			ctx:       context.Background(),
			setupMock: func(repo *profilemocks.ProfileRepository) {},
			assertFn: func(t *testing.T, result profile.ProfileResponse, err error) {
				require.Error(t, err)
				testingutil.AssertAppError(t, err, apperror.CodeForbidden, http.StatusForbidden, "unauthorized")
			},
		},
		{
			name: "not found",
			ctx:  ctxWithUser(userID),
			setupMock: func(repo *profilemocks.ProfileRepository) {
				repo.On("GetByID", mock.Anything, userID).Return(nil, profile.ErrProfileNotFound).Once()
			},
			assertFn: func(t *testing.T, result profile.ProfileResponse, err error) {
				require.Error(t, err)
				testingutil.AssertAppError(t, err, apperror.CodeNotFound, http.StatusNotFound, "profile not found")
			},
		},
		{
			name: "repository error",
			ctx:  ctxWithUser(userID),
			setupMock: func(repo *profilemocks.ProfileRepository) {
				repo.On("GetByID", mock.Anything, userID).Return(nil, errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, result profile.ProfileResponse, err error) {
				require.Error(t, err)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := profilemocks.NewProfileRepository(t)
			storage := storagemocks.NewStorage(t)
			svc := newProfileService(repo, storage)
			tc.setupMock(repo)

			result, err := svc.GetProfile(tc.ctx)
			tc.assertFn(t, result, err)
		})
	}
}

func TestService_UpdateProfile(t *testing.T) {
	userID := uuid.New()
	phone := "+1234567890"

	testCases := []struct {
		name      string
		ctx       context.Context
		req       *profile.UpdateProfileRequest
		setupMock func(*profilemocks.ProfileRepository)
		assertFn  func(*testing.T, error)
	}{
		{
			name: "success",
			ctx:  ctxWithUser(userID),
			req:  &profile.UpdateProfileRequest{Name: "Alice", Phone: &phone},
			setupMock: func(repo *profilemocks.ProfileRepository) {
				repo.On("UpdateProfile", mock.Anything, userID, "Alice", &phone).Return(nil).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:      "no user in context",
			ctx:       context.Background(),
			req:       &profile.UpdateProfileRequest{Name: "Alice"},
			setupMock: func(repo *profilemocks.ProfileRepository) {},
			assertFn: func(t *testing.T, err error) {
				require.Error(t, err)
				testingutil.AssertAppError(t, err, apperror.CodeForbidden, http.StatusForbidden, "unauthorized")
			},
		},
		{
			name: "not found",
			ctx:  ctxWithUser(userID),
			req:  &profile.UpdateProfileRequest{Name: "Alice"},
			setupMock: func(repo *profilemocks.ProfileRepository) {
				repo.On("UpdateProfile", mock.Anything, userID, "Alice", mock.Anything).Return(profile.ErrProfileNotFound).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.Error(t, err)
				testingutil.AssertAppError(t, err, apperror.CodeNotFound, http.StatusNotFound, "profile not found")
			},
		},
		{
			name: "repository error",
			ctx:  ctxWithUser(userID),
			req:  &profile.UpdateProfileRequest{Name: "Alice"},
			setupMock: func(repo *profilemocks.ProfileRepository) {
				repo.On("UpdateProfile", mock.Anything, userID, "Alice", mock.Anything).Return(errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := profilemocks.NewProfileRepository(t)
			storage := storagemocks.NewStorage(t)
			svc := newProfileService(repo, storage)
			tc.setupMock(repo)

			err := svc.UpdateProfile(tc.ctx, tc.req)
			tc.assertFn(t, err)
		})
	}
}
