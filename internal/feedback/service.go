package feedback

import (
	"context"
	"errors"

	"github.com/riyanamanda/helpdesk-backend/internal/notification"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
)

type FeedbackService interface {
	ListFeedbacks(ctx context.Context, params *GetFeedbackParams) ([]FeedbackResponse, int64, error)
	CreateFeedback(ctx context.Context, req *CreateFeedbackRequest) error
	GetFeedback(ctx context.Context, id int64) (*FeedbackResponse, error)
	UpdateFeedbackStatus(ctx context.Context, id int64, req UpdateFeedbackStatusRequest) error
	DeleteFeedback(ctx context.Context, id int64) error
}

type service struct {
	repo            FeedbackRepository
	notificationSvc notification.Notifier
}

func NewFeedbackService(repo FeedbackRepository, notificationSvc notification.Notifier) FeedbackService {
	return &service{repo: repo, notificationSvc: notificationSvc}
}

func (s *service) ListFeedbacks(ctx context.Context, params *GetFeedbackParams) ([]FeedbackResponse, int64, error) {
	if params == nil {
		params = &GetFeedbackParams{}
	}
	params.Normalize()

	feedbacks, total, err := s.repo.GetAll(ctx, *params)
	if err != nil {
		return nil, 0, err
	}

	return toFeedbackResponses(feedbacks), total, nil
}

func (s *service) CreateFeedback(ctx context.Context, req *CreateFeedbackRequest) error {
	createdBy, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	feedback := Feedback{
		Title:       req.Title,
		Description: req.Description,
		Type:        FeedbackType(req.Type),
		CreatedBy:   createdBy,
	}

	return s.repo.Create(ctx, feedback)
}

func (s *service) GetFeedback(ctx context.Context, id int64) (*FeedbackResponse, error) {
	feedback, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrFeedbackNotFound) {
			return nil, apperr.NotFound("feedback")
		}
		return nil, err
	}

	result := toFeedbackResponse(*feedback)
	return &result, nil
}

func (s *service) UpdateFeedbackStatus(ctx context.Context, id int64, req UpdateFeedbackStatusRequest) error {
	reviewerID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrFeedbackNotFound) {
			return apperr.NotFound("feedback")
		}
		return err
	}

	if FeedbackStatus(existing.Status) == FeedbackStatusDelivered || FeedbackStatus(existing.Status) == FeedbackStatusRejected {
		return apperr.BadRequest("status is delivered or rejected")
	}

	if err := s.repo.UpdateStatus(ctx, id, reviewerID, req.Status); err != nil {
		if errors.Is(err, ErrFeedbackNotFound) {
			return apperr.NotFound("feedback")
		}
		return err
	}

	s.notificationSvc.FeedbackStatusUpdated(ctx, id, existing.CreatedByID, reviewerID, string(req.Status))
	return nil
}

func (s *service) DeleteFeedback(ctx context.Context, id int64) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrFeedbackNotFound) {
			return apperr.NotFound("feedback")
		}
		return err
	}

	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	if existing.CreatedByID != userID {
		return apperr.Forbidden("you can only delete your own feedback")
	}

	if FeedbackStatus(existing.Status) != FeedbackStatusOpen {
		return apperr.BadRequest("only open feedback can be deleted")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
