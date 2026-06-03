package profile

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperror"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/firebase"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/upload"
	"github.com/riyanamanda/helpdesk-backend/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type ProfileService interface {
	GetProfile(ctx context.Context) (ProfileResponse, error)
	UpdateProfile(ctx context.Context, req *UpdateProfileRequest) error
	UpdateAvatar(ctx context.Context, file *upload.File) error
	SyncGoogle(ctx context.Context, req *SyncGoogleRequest) error
	RevokeGoogle(ctx context.Context) error
	UpdatePassword(ctx context.Context, req UpdatePasswordRequest) error
}

type service struct {
	repo          ProfileRepository
	storage       storage.Storage
	storageConfig config.Storage
	authConfig    config.Auth
}

func NewProfileService(repo ProfileRepository, store storage.Storage, storageConfig config.Storage, authConfig config.Auth) ProfileService {
	return &service{
		repo:          repo,
		storage:       store,
		storageConfig: storageConfig,
		authConfig:    authConfig,
	}
}

func (s *service) GetProfile(ctx context.Context) (ProfileResponse, error) {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return ProfileResponse{}, apperror.Unauthorized(apperror.CodeUnauthorized, "unauthorized")
	}

	p, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			return ProfileResponse{}, apperror.NotFound("profile")
		}
		return ProfileResponse{}, err
	}

	return toProfileResponse(*p, s.storageConfig), nil
}

func (s *service) UpdateProfile(ctx context.Context, req *UpdateProfileRequest) error {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperror.Unauthorized(apperror.CodeUnauthorized, "unauthorized")
	}

	currentUser, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if req.Email != currentUser.Email {
		if currentUser.GoogleID != nil {
			return apperror.Forbidden("Please unlink your google before change your email")
		}
	}

	if err := s.repo.UpdateProfile(ctx, userID, req.Name, req.Email, req.Phone, strings.ToUpper(req.Gender)); err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			return apperror.NotFound("profile")
		}
		return err
	}

	return nil
}

func (s *service) UpdateAvatar(ctx context.Context, file *upload.File) error {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperror.Unauthorized(apperror.CodeUnauthorized, "unauthorized")
	}

	objectKey := fmt.Sprintf("avatars/%s/avatar", userID.String())
	if err := s.storage.Upload(ctx, objectKey, file.Content, file.Size, file.ContentType); err != nil {
		return err
	}

	return s.repo.UpdateAvatar(ctx, userID, objectKey)
}

func (s *service) SyncGoogle(ctx context.Context, req *SyncGoogleRequest) error {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperror.Unauthorized(apperror.CodeUnauthorized, "unauthorized")
	}

	claims, err := firebase.VerifyIDToken(req.IDToken, s.authConfig.FirebaseProjectID)
	if err != nil {
		return apperror.Unauthorized(apperror.CodeUnauthorized, "invalid google token")
	}

	currentProfile, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if currentProfile.Email != claims.Email {
		return apperror.BadRequest("google account email does not match your account email")
	}

	if err := s.repo.SetGoogleID(ctx, userID, claims.Subject); err != nil {
		if errors.Is(err, ErrGoogleIDAlreadyLinked) {
			return apperror.AlreadyExists("google account")
		}
		return err
	}

	return nil
}

func (s *service) RevokeGoogle(ctx context.Context) error {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperror.Unauthorized(apperror.CodeUnauthorized, "unauthorized")
	}

	currentProfile, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if currentProfile.GoogleID == nil {
		return apperror.NotFound("your account is not linked to google")
	}

	if err := s.repo.UnsetGoogleID(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (s *service) UpdatePassword(ctx context.Context, req UpdatePasswordRequest) error {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperror.Unauthorized(apperror.CodeUnauthorized, "unauthorized")
	}

	currentUser, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(req.CurrentPassword)); err != nil {
		return apperror.BadRequest("invalid password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return err
	}

	return nil
}
