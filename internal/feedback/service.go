package feedback

import (
	"context"
	"errors"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
)

type FeedbackService interface {
	ListFeedbacks(ctx context.Context, params *GetFeedbackParams) ([]FeedbackResponse, int64, error)
	CreateFeedback(ctx context.Context, req *CreateFeedbackRequest) error
	GetFeedback(ctx context.Context, id int64) (*FeedbackResponse, error)
	UpdateFeedbackStatus(ctx context.Context, id int64, req UpdateFeedbackStatusRequest) error
}

type service struct {
	repo FeedbackRepository
}

func NewFeedbackService(repo FeedbackRepository) FeedbackService {
	return &service{repo: repo}
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

	if err := s.repo.UpdateStatus(ctx, id, reviewerID, req.Status); err != nil {
		if errors.Is(err, ErrFeedbackNotFound) {
			return apperr.NotFound("feedback")
		}
		return err
	}

	return nil
}
