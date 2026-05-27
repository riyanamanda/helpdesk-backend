package profile

import (
	"context"
	"errors"
	"fmt"

	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/upload"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
	"github.com/riyanamanda/helpdesk-backend/internal/storage"
)

type ProfileService interface {
	GetProfile(ctx context.Context) (ProfileResponse, error)
	UpdateProfile(ctx context.Context, req *UpdateProfileRequest) error
	UpdateAvatar(ctx context.Context, file *upload.File) error
	SyncGoogle(ctx context.Context, req *SyncGoogleRequest) error
}

type service struct {
	profileRepo   ProfileRepository
	storage       storage.Storage
	storageConfig config.Storage
	authConfig    config.Auth
}

func NewProfileService(
	profileRepo ProfileRepository,
	storage storage.Storage,
	storageConfig config.Storage,
	authConfig config.Auth,
) ProfileService {
	return &service{
		profileRepo:   profileRepo,
		storage:       storage,
		storageConfig: storageConfig,
		authConfig:    authConfig,
	}
}

func (s *service) GetProfile(ctx context.Context) (ProfileResponse, error) {
	userID, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		return ProfileResponse{}, apperror.Forbidden("unauthorized")
	}

	p, err := s.profileRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			return ProfileResponse{}, apperror.NotFound("profile")
		}
		return ProfileResponse{}, err
	}

	return toProfileResponse(*p, s.storageConfig), nil
}

func (s *service) UpdateProfile(ctx context.Context, req *UpdateProfileRequest) error {
	userID, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		return apperror.Forbidden("unauthorized")
	}

	if err := s.profileRepo.UpdateProfile(ctx, userID, req.Name, req.Phone); err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			return apperror.NotFound("profile")
		}
		return err
	}

	return nil
}

func (s *service) UpdateAvatar(ctx context.Context, file *upload.File) error {
	userID, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		return apperror.Forbidden("unauthorized")
	}

	objectKey := fmt.Sprintf("avatars/%s/avatar", userID.String())
	if err := s.storage.Upload(ctx, objectKey, file.Content, file.Size, file.ContentType); err != nil {
		return err
	}

	return s.profileRepo.UpdateAvatar(ctx, userID, objectKey)
}

func (s *service) SyncGoogle(ctx context.Context, req *SyncGoogleRequest) error {
	userID, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		return apperror.Forbidden("unauthorized")
	}

	claims, err := utils.VerifyFirebaseIDToken(req.IDToken, s.authConfig.FirebaseProjectID)
	if err != nil {
		return apperror.Unauthorized(apperror.CodeUnauthorized, "invalid google token")
	}

	currentProfile, err := s.profileRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if currentProfile.Email != claims.Email {
		return apperror.BadRequest("google account email does not match your account email")
	}

	if err := s.profileRepo.SetGoogleID(ctx, userID, claims.Subject); err != nil {
		if errors.Is(err, ErrGoogleIDAlreadyLinked) {
			return apperror.AlreadyExists("google account")
		}
		return err
	}

	return nil
}
