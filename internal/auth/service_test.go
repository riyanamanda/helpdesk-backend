package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	auth "github.com/riyanamanda/helpdesk-backend/internal/auth"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	testingutil "github.com/riyanamanda/helpdesk-backend/internal/shared/testing"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
	user "github.com/riyanamanda/helpdesk-backend/internal/user"
	usermocks "github.com/riyanamanda/helpdesk-backend/internal/user/mocks"
)

func TestService_Login(t *testing.T) {
	secret := "test-secret"
	expiresIn := 15 * time.Minute
	userID := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	t.Run("success with configured jwt secret and expiry", func(t *testing.T) {
		repo := usermocks.NewUserRepository(t)
		repo.On("GetByEmail", mock.Anything, "admin@email.com").Return(&user.User{
			ID:       userID,
			Email:    "admin@email.com",
			Password: string(hashedPassword),
			Role:     user.ADMIN,
			IsActive: true,
		}, nil).Once()

		svc := auth.NewAuthService(repo, secret, expiresIn)

		result, err := svc.Login(context.Background(), &auth.LoginRequest{
			Email:    "admin@email.com",
			Password: "password123",
		})

		require.NoError(t, err)
		require.NotEmpty(t, result.AccessToken)

		claims, err := utils.ParseToken(result.AccessToken, secret)
		require.NoError(t, err)
		assert.Equal(t, userID.String(), claims.Subject)
		assert.Equal(t, string(user.ADMIN), claims.Role)
		require.NotNil(t, claims.ExpiresAt)
		assert.WithinDuration(t, time.Now().Add(expiresIn), claims.ExpiresAt.Time, 2*time.Second)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		repo := usermocks.NewUserRepository(t)
		repo.On("GetByEmail", mock.Anything, "missing@email.com").Return(nil, user.ErrUserNotFound).Once()

		svc := auth.NewAuthService(repo, secret, expiresIn)

		result, err := svc.Login(context.Background(), &auth.LoginRequest{
			Email:    "missing@email.com",
			Password: "password123",
		})

		require.Error(t, err)
		assert.Equal(t, auth.LoginResponse{}, result)
		testingutil.AssertAppError(t, err, apperror.CODE_BAD_REQUEST, http.StatusBadRequest, "invalid email or password")
	})
}
