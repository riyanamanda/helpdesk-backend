package profile

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/firebase"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/jwtutil"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*user.UserProjection, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, name string, email string, phone *string, gender string) error
	UpdateAvatar(ctx context.Context, id uuid.UUID, avatarKey string) error
	SetGoogleID(ctx context.Context, id uuid.UUID, googleID string) error
	UnsetGoogleID(ctx context.Context, id uuid.UUID) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, password string) error
}

type Service struct {
	repo          repository
	storage       storage.Storage
	storageConfig config.Storage
	authConfig    config.Auth
	cache         cache.Cache
}

func NewService(repo repository, store storage.Storage, storageConfig config.Storage, authConfig config.Auth, cache cache.Cache) *Service {
	return &Service{
		repo:          repo,
		storage:       store,
		storageConfig: storageConfig,
		authConfig:    authConfig,
		cache:         cache,
	}
}

func (s *Service) GetProfile(ctx context.Context) (*ProfileResponse, error) {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return nil, apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	p, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			return nil, apperr.NotFound("profile")
		}
		return nil, err
	}

	result := toProfileResponse(*p, s.storageConfig)

	return &result, nil
}

func (s *Service) UpdateProfile(ctx context.Context, req *UpdateProfileRequest) error {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	currentUser, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if req.Email != currentUser.Email {
		if currentUser.GoogleID != nil {
			return apperr.Forbidden("Please unlink your google before change your email")
		}
	}

	if err := s.repo.UpdateProfile(ctx, userID, req.Name, req.Email, req.Phone, strings.ToUpper(req.Gender)); err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			return apperr.NotFound("profile")
		}
		return err
	}

	user.InvalidateCache(ctx, s.cache)

	return nil
}

func (s *Service) UpdateAvatar(ctx context.Context, file *storage.File) error {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	objectKey := fmt.Sprintf("avatars/%s/avatar", userID.String())
	if err := s.storage.Upload(ctx, objectKey, file); err != nil {
		return err
	}

	if err := s.repo.UpdateAvatar(ctx, userID, objectKey); err != nil {
		return err
	}

	user.InvalidateCache(ctx, s.cache)

	return nil
}

func (s *Service) SyncGoogle(ctx context.Context, req *SyncGoogleRequest) error {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	claims, err := firebase.VerifyIDToken(req.IDToken, s.authConfig.FirebaseProjectID)
	if err != nil {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "invalid google token")
	}

	currentProfile, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if currentProfile.Email != claims.Email {
		return apperr.BadRequest("google account email does not match your account email")
	}

	if err := s.repo.SetGoogleID(ctx, userID, claims.Subject); err != nil {
		if errors.Is(err, ErrGoogleIDAlreadyLinked) {
			return apperr.AlreadyExists("google account")
		}
		return err
	}

	return nil
}

func (s *Service) RevokeGoogle(ctx context.Context) error {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	currentProfile, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if currentProfile.GoogleID == nil {
		return apperr.BadRequest("your account is not linked to google")
	}

	return s.repo.UnsetGoogleID(ctx, userID)
}

func (s *Service) UpdatePassword(ctx context.Context, req UpdatePasswordRequest) error {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	currentUser, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(req.CurrentPassword)); err != nil {
		return apperr.BadRequest("invalid password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return err
	}

	return jwtutil.RevokeUserSessions(ctx, s.cache, userID)
}
