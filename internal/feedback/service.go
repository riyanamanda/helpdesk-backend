package feedback

import (
	"context"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/event"
	"github.com/riyanamanda/helpdesk-backend/internal/outbox"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
)

type repository interface {
	GetAll(ctx context.Context, params GetFeedbackParams) ([]FeedbackProjection, int64, error)
	GetByID(ctx context.Context, id int64) (*FeedbackProjection, error)
	Create(ctx context.Context, tx database.Tx, feedback Feedback) (int64, error)
	UpdateStatus(ctx context.Context, tx database.Tx, id int64, reviewerID uuid.UUID, status FeedbackStatus) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	repo       repository
	outboxRepo outbox.Repository
	txManager  *database.TxManager
}

func NewService(repo repository, outboxRepo outbox.Repository, txManager *database.TxManager) *Service {
	return &Service{
		repo:       repo,
		outboxRepo: outboxRepo,
		txManager:  txManager,
	}
}

func (s *Service) addOutboxEvent(ctx context.Context, tx database.Tx, eventType string, feedbackID int64, payload any) error {
	ev, err := outbox.NewEvent(eventType, feedbackID, payload)
	if err != nil {
		return err
	}

	return s.outboxRepo.Create(ctx, tx, ev)
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

	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	id, err := s.repo.Create(ctx, tx, feedback)
	if err != nil {
		return err
	}

	if err := s.addOutboxEvent(ctx, tx, event.FeedbackCreated, id, event.FeedbackCreatedEvent{
		FeedbackID:    id,
		Title:         feedback.Title,
		Type:          string(feedback.Type),
		Description:   feedback.Description,
		SubmittedByID: createdBy,
	}); err != nil {
		return err
	}

	return tx.Commit()
}

func canAccess(ctx context.Context, feedback *FeedbackProjection) bool {
	authUser, ok := ctxkey.GetAuthUserFromContext(ctx)
	if !ok || authUser == nil {
		return false
	}

	return feedback.CreatedByID == authUser.ID || authUser.Permissions.Has(rbac.PermissionFeedbackUpdate)
}

func (s *Service) GetFeedback(ctx context.Context, id int64) (*FeedbackResponse, error) {
	feedback, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !canAccess(ctx, feedback) {
		return nil, ErrFeedbackNotFound
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
		return err
	}

	if FeedbackStatus(existing.Status) == FeedbackStatusDelivered || FeedbackStatus(existing.Status) == FeedbackStatusRejected {
		return apperr.BadRequest("status is delivered or rejected")
	}

	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.repo.UpdateStatus(ctx, tx, id, reviewerID, req.Status); err != nil {
		return err
	}

	if existing.CreatedByID != reviewerID {
		if err := s.addOutboxEvent(ctx, tx, event.FeedbackStatusChanged, id, event.FeedbackStatusChangedEvent{
			FeedbackID:   id,
			Title:        existing.Title,
			Status:       string(req.Status),
			SubmitterID:  existing.CreatedByID,
			ReviewedByID: reviewerID,
		}); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Service) DeleteFeedback(ctx context.Context, id int64) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !canAccess(ctx, existing) {
		return ErrFeedbackNotFound
	}

	if FeedbackStatus(existing.Status) != FeedbackStatusOpen {
		return apperr.BadRequest("only open feedback can be deleted")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
