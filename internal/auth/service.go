package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/riyanamanda/helpdesk-backend/internal/event"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/firebase"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/rabbitmq"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/jwtutil"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type AuthService interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	LoginWithGoogle(ctx context.Context, req *GoogleLoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context) error
	Me(ctx context.Context) (*CurrentUserResponse, error)
	ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req ResetPasswordRequest) error
}

type service struct {
	userRepo          user.UserRepository
	authConfig        config.Auth
	storageConfig     config.Storage
	appConfig         config.App
	redis             cache.Cache
	permissionService ctxkey.PermissionService
	rabbitmq          rabbitmq.Client
}

func NewAuthService(
	repo user.UserRepository,
	authConfig config.Auth,
	storageConfig config.Storage,
	appConfig config.App,
	redis cache.Cache,
	permissionService ctxkey.PermissionService,
	rabbitmq rabbitmq.Client,
) AuthService {
	return &service{
		userRepo:          repo,
		authConfig:        authConfig,
		storageConfig:     storageConfig,
		appConfig:         appConfig,
		redis:             redis,
		permissionService: permissionService,
		rabbitmq:          rabbitmq,
	}
}

func (s *service) issueSession(ctx context.Context, user user.UserProjection) (string, error) {
	if !user.IsActive {
		return "", apperr.Forbidden("user is inactive")
	}

	token, jti, err := jwtutil.GenerateToken(user.ID, s.authConfig.JWTSecret, s.authConfig.JWTExp)
	if err != nil {
		return "", err
	}

	if err := s.redis.Set(ctx, jwtutil.TokenKeyPrefix+jti, user.ID.String(), s.authConfig.JWTExp); err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	currentUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, apperr.BadRequest("invalid email or password")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(req.Password)); err != nil {
		return nil, apperr.BadRequest("invalid email or password")
	}

	token, err := s.issueSession(ctx, *currentUser)
	if err != nil {
		return nil, err
	}

	permissions, err := s.permissionService.GetUserPermissions(ctx, currentUser.ID)
	if err != nil {
		return nil, err
	}

	return toLoginResponse(token, *currentUser, s.storageConfig, permissions.ToSlice()), nil
}

func (s *service) LoginWithGoogle(ctx context.Context, req *GoogleLoginRequest) (*LoginResponse, error) {
	firebaseClaims, err := firebase.VerifyIDToken(req.IDToken, s.authConfig.FirebaseProjectID)
	if err != nil {
		return nil, apperr.Unauthorized(apperr.CodeUnauthorized, "invalid google token")
	}

	currentUser, err := s.userRepo.GetByEmail(ctx, firebaseClaims.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, apperr.Forbidden("your google account is not registered")
		}
		return nil, err
	}

	if currentUser.GoogleID == nil {
		return nil, apperr.Forbidden("google account is not linked to this account")
	}

	token, err := s.issueSession(ctx, *currentUser)
	if err != nil {
		return nil, err
	}

	permissions, err := s.permissionService.GetUserPermissions(ctx, currentUser.ID)
	if err != nil {
		return nil, err
	}

	return toLoginResponse(token, *currentUser, s.storageConfig, permissions.ToSlice()), nil
}

func (s *service) Logout(ctx context.Context) error {
	jti, ok := ctxkey.GetJTIFromContext(ctx)
	if !ok || jti == "" {
		return apperr.Unauthorized(apperr.CodeInvalidToken, "invalid token")
	}

	return s.redis.Delete(ctx, jwtutil.TokenKeyPrefix+jti)
}

func (s *service) Me(ctx context.Context) (*CurrentUserResponse, error) {
	authUser, ok := ctxkey.GetAuthUserFromContext(ctx)
	if !ok || authUser == nil {
		return nil, apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	u, err := s.userRepo.GetByID(ctx, authUser.ID)
	if err != nil {
		return nil, err
	}

	return toCurrentUserResponse(*u, s.storageConfig, authUser.Permissions.ToSlice()), nil
}

func (s *service) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {
	userValue, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// prevent user enumeration
		if errors.Is(err, user.ErrUserNotFound) {
			return nil
		}

		return err
	}

	// invalidate old token if exists
	userKey := buildPasswordResetRequestCache(userValue.ID.String())
	existingToken, err := s.redis.Get(ctx, userKey)
	if err == nil && existingToken != "" {
		InvalidateResetPasswordCache(ctx, s.redis, userValue.ID.String(), existingToken)
	}

	// generate new reset password token
	tokenByte := make([]byte, 32)
	if _, err := rand.Read(tokenByte); err != nil {
		return err
	}
	token := hex.EncodeToString(tokenByte)

	// set request token
	if err := s.redis.Set(ctx, buildPasswordResetRequestCache(userValue.ID.String()), token, 60*time.Second); err != nil {
		return err
	}

	// set reset token
	if err := s.redis.Set(ctx, buildPasswordResetCache(token), userValue.ID.String(), 15*time.Minute); err != nil {
		_ = s.redis.Delete(ctx, buildPasswordResetRequestCache(userValue.ID.String()))
		return err
	}

	resetUrl := fmt.Sprintf("%s/reset-password?token=%s", s.appConfig.URL, url.QueryEscape(token))
	passwordEvent := event.PasswordResetRequestedEvent{
		Name:     userValue.Name,
		Email:    userValue.Email,
		ResetURL: resetUrl,
	}
	payload, err := json.Marshal(passwordEvent)
	if err != nil {
		return err
	}

	if err := s.rabbitmq.Publish(ctx, rabbitmq.ExchangeEvent, event.PasswordResetRequested, "application/json", payload); err != nil {
		return err
	}

	return nil
}

func (s *service) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	return nil
}
