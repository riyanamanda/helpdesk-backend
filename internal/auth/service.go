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

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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

type repository interface {
	GetByEmail(ctx context.Context, email string) (*user.UserProjection, error)
	GetByID(ctx context.Context, id uuid.UUID) (*user.UserProjection, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
}

type Service struct {
	userRepo          repository
	authConfig        config.Auth
	storageConfig     config.Storage
	appConfig         config.App
	redis             cache.Cache
	permissionService ctxkey.PermissionService
	rabbitmq          rabbitmq.Client
}

func NewService(
	repo repository,
	authConfig config.Auth,
	storageConfig config.Storage,
	appConfig config.App,
	redis cache.Cache,
	permissionService ctxkey.PermissionService,
	rabbitmq rabbitmq.Client,
) *Service {
	return &Service{
		userRepo:          repo,
		authConfig:        authConfig,
		storageConfig:     storageConfig,
		appConfig:         appConfig,
		redis:             redis,
		permissionService: permissionService,
		rabbitmq:          rabbitmq,
	}
}

func (s *Service) issueSession(ctx context.Context, user user.UserProjection) (string, error) {
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

	if err := jwtutil.TrackSession(ctx, s.redis, user.ID, jti, s.authConfig.JWTExp); err != nil {
		_ = s.redis.Delete(ctx, jwtutil.TokenKeyPrefix+jti)
		return "", err
	}

	return token, nil
}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
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

func (s *Service) LoginWithGoogle(ctx context.Context, req *GoogleLoginRequest) (*LoginResponse, error) {
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

func (s *Service) Logout(ctx context.Context) error {
	jti, ok := ctxkey.GetJTIFromContext(ctx)
	if !ok || jti == "" {
		return apperr.Unauthorized(apperr.CodeInvalidToken, "invalid token")
	}

	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	if err := s.redis.Delete(ctx, jwtutil.TokenKeyPrefix+jti); err != nil {
		return err
	}

	return jwtutil.UntrackSession(ctx, s.redis, userID, jti)
}

func (s *Service) Me(ctx context.Context) (*CurrentUserResponse, error) {
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

func (s *Service) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {
	userValue, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil
		}

		return err
	}

	// Check cooldown
	cooldownKey := buildPasswordResetCooldownCache(userValue.ID.String())

	ttl, err := s.redis.TTL(ctx, cooldownKey)
	if err != nil {
		return err
	}

	if ttl > 0 {
		return apperr.RateLimited(
			"your request is in cooldown",
			int64(ttl.Seconds())+1,
		)
	}

	// Check active token
	requestKey := buildPasswordResetRequestCache(userValue.ID.String())

	existingToken, err := s.redis.Get(ctx, requestKey)
	if err == nil && existingToken != "" {
		return apperr.BadRequest("you already have an active reset token")
	}

	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	// Generate token
	tokenByte := make([]byte, 32)
	if _, err := rand.Read(tokenByte); err != nil {
		return err
	}

	token := hex.EncodeToString(tokenByte)

	// Active token: 15 minutes
	if err := s.redis.Set(ctx, requestKey, token, 15*time.Minute); err != nil {
		return err
	}

	if err := s.redis.Set(ctx, buildPasswordResetTokenCache(token), userValue.ID.String(), 15*time.Minute); err != nil {
		_ = s.redis.Delete(ctx, requestKey)
		return err
	}

	// Cooldown: 60 seconds
	if err := s.redis.Set(ctx, cooldownKey, "1", 60*time.Second); err != nil {
		_ = s.redis.DeleteMany(ctx, requestKey, buildPasswordResetTokenCache(token))
		return err
	}

	resetURL := fmt.Sprintf(
		"%s/reset-password?token=%s",
		s.appConfig.URL,
		url.QueryEscape(token),
	)

	passwordEvent := event.PasswordResetRequestedEvent{
		Name:     userValue.Name,
		Email:    userValue.Email,
		ResetURL: resetURL,
	}

	payload, err := json.Marshal(passwordEvent)
	if err != nil {
		return err
	}

	if err := s.rabbitmq.Publish(
		ctx,
		rabbitmq.ExchangeEvent,
		event.PasswordResetRequested,
		"application/json",
		payload,
	); err != nil {
		InvalidateResetPasswordCache(ctx, s.redis, userValue.ID.String(), token)
		return err
	}

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	userID, err := s.redis.Get(ctx, buildPasswordResetTokenCache(req.Token))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return apperr.TokenExpired("token invalid or expired, please make a reset request again")
		}

		return err
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return apperr.TokenExpired("token invalid or expired, please make a reset request again")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(ctx, id, string(hashedPassword)); err != nil {
		return err
	}

	if err := jwtutil.RevokeUserSessions(ctx, s.redis, id); err != nil {
		return err
	}

	InvalidateResetPasswordCache(ctx, s.redis, id.String(), req.Token)

	return nil
}
