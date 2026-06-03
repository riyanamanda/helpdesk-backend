package auth

import (
	"context"
	"errors"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/firebase"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/jwtutil"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

const tokenKeyPrefix = "auth:token:"

type sessionStore interface {
	Set(ctx context.Context, key, value string, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

type redisAdapter struct {
	client *goredis.Client
}

func (r *redisAdapter) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisAdapter) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

type AuthService interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	LoginWithGoogle(ctx context.Context, req *GoogleLoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context) error
	Me(ctx context.Context) (*CurrentUserResponse, error)
}

type service struct {
	userRepo      user.UserRepository
	config        config.Auth
	storageConfig config.Storage
	redis         sessionStore
}

func NewAuthService(repo user.UserRepository, cfg config.Auth, storageConfig config.Storage, redis sessionStore) AuthService {
	return &service{
		userRepo:      repo,
		config:        cfg,
		storageConfig: storageConfig,
		redis:         redis,
	}
}

func (s *service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	currentUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, apperr.BadRequest("invalid email or password")
		}
		return nil, err
	}

	if !currentUser.IsActive {
		return nil, apperr.Forbidden("user is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(req.Password)); err != nil {
		return nil, apperr.BadRequest("invalid email or password")
	}

	token, jti, err := jwtutil.GenerateToken(currentUser.ID, string(currentUser.Role), s.config.JWTSecret, s.config.JWTExp)
	if err != nil {
		return nil, err
	}

	if err := s.redis.Set(ctx, tokenKeyPrefix+jti, currentUser.ID.String(), s.config.JWTExp); err != nil {
		return nil, err
	}

	result := LoginResponse{
		User:        toCurrentUserResponse(*currentUser, s.storageConfig),
		AccessToken: token,
	}

	return &result, nil
}

func (s *service) LoginWithGoogle(ctx context.Context, req *GoogleLoginRequest) (*LoginResponse, error) {
	firebaseClaims, err := firebase.VerifyIDToken(req.IDToken, s.config.FirebaseProjectID)
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

	if !currentUser.IsActive {
		return nil, apperr.Forbidden("user is inactive")
	}

	token, jti, err := jwtutil.GenerateToken(currentUser.ID, string(currentUser.Role), s.config.JWTSecret, s.config.JWTExp)
	if err != nil {
		return nil, err
	}

	if err := s.redis.Set(ctx, tokenKeyPrefix+jti, currentUser.ID.String(), s.config.JWTExp); err != nil {
		return nil, err
	}

	result := LoginResponse{
		User:        toCurrentUserResponse(*currentUser, s.storageConfig),
		AccessToken: token,
	}

	return &result, nil
}

func (s *service) Logout(ctx context.Context) error {
	jti, ok := ctxkey.GetJTIFromContext(ctx)
	if !ok || jti == "" {
		return apperr.Unauthorized(apperr.CodeInvalidToken, "invalid token")
	}

	return s.redis.Delete(ctx, tokenKeyPrefix+jti)
}

func (s *service) Me(ctx context.Context) (*CurrentUserResponse, error) {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return nil, apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := toCurrentUserResponse(*u, s.storageConfig)

	return &result, nil
}
