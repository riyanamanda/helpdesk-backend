package feedback

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
)

type repository interface {
	GetAll(ctx context.Context, params GetFeedbackParams) ([]FeedbackProjection, int64, error)
	GetByID(ctx context.Context, id int64) (*FeedbackProjection, error)
	Create(ctx context.Context, feedback Feedback) error
	UpdateStatus(ctx context.Context, id int64, reviewerID uuid.UUID, status FeedbackStatus) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListFeedbacks(ctx context.Context, params *GetFeedbackParams) ([]FeedbackResponse, int64, error) {
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

func (s *Service) CreateFeedback(ctx context.Context, req *CreateFeedbackRequest) error {
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

func (s *Service) GetFeedback(ctx context.Context, id int64) (*FeedbackResponse, error) {
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

func (s *Service) UpdateFeedbackStatus(ctx context.Context, id int64, req UpdateFeedbackStatusRequest) error {
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

	return nil
}

func (s *Service) DeleteFeedback(ctx context.Context, id int64) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrFeedbackNotFound) {
			return apperr.NotFound("feedback")
		}
		return err
	}

	if FeedbackStatus(existing.Status) != FeedbackStatusOpen {
		return apperr.BadRequest("only open feedback can be deleted")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
