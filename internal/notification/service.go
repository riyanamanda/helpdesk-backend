package notification

import (
	"context"
	"errors"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
)

type NotificationService interface {
	ListNotifications(ctx context.Context) ([]NotificationResponse, error)
	UnreadCount(ctx context.Context) (UnreadCountResponse, error)
	MarkAsRead(ctx context.Context, id int64) error
	MarkAllAsRead(ctx context.Context) error
}

type service struct {
	repo NotificationRepository
}

func NewNotificationService(repo NotificationRepository) NotificationService {
	return &service{
		repo: repo,
	}
}

func (s *service) UnreadCount(ctx context.Context) (UnreadCountResponse, error) {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return UnreadCountResponse{}, apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	count, err := s.repo.CountUnread(ctx, userID)
	if err != nil {
		return UnreadCountResponse{}, err
	}

	return UnreadCountResponse{Count: count}, nil
}

func (s *service) MarkAllAsRead(ctx context.Context) error {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	return s.repo.MarkAllAsRead(ctx, userID)
}

func (s *service) MarkAsRead(ctx context.Context, id int64) error {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	if err := s.repo.MarkAsRead(ctx, id, userID); err != nil {
		if errors.Is(err, ErrNotificationNotFound) {
			return apperr.NotFound("notification")
		}
		return err
	}

	return nil
}

func (s *service) ListNotifications(ctx context.Context) ([]NotificationResponse, error) {
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return nil, apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	notifications, err := s.repo.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	return toNotificationResponses(notifications), nil
}
