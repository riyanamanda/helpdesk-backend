package user_device

import (
	"context"
	"errors"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
)

type UserDeviceService interface {
	RegisterDevice(ctx context.Context, req RegisterDeviceRequest) error
	UnregisterDevice(ctx context.Context, req UnregisterDeviceRequest) error
}

type service struct {
	repo UserDeviceRepository
}

func NewUserDeviceService(repo UserDeviceRepository) UserDeviceService {
	return &service{repo: repo}
}

func (s *service) RegisterDevice(ctx context.Context, req RegisterDeviceRequest) error {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}
	return s.repo.Upsert(ctx, userID, req.FcmToken)
}

func (s *service) UnregisterDevice(ctx context.Context, req UnregisterDeviceRequest) error {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	if err := s.repo.Delete(ctx, userID, req.FcmToken); err != nil {
		if errors.Is(err, ErrDeviceNotFound) {
			return apperr.NotFound("device")
		}
		return err
	}
	return nil
}
