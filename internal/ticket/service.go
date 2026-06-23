package ticket

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/riyanamanda/helpdesk-backend/internal/category"
	"github.com/riyanamanda/helpdesk-backend/internal/dashboard"
	"github.com/riyanamanda/helpdesk-backend/internal/division"
	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
	"github.com/riyanamanda/helpdesk-backend/internal/notification"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type TicketService interface {
	ListTickets(ctx context.Context, params *GetTicketParams) ([]TicketResponse, int64, error)
	CreateTicket(ctx context.Context, req *TicketCreateRequest, file *storage.File) error
	GetTicket(ctx context.Context, id int64) (*TicketDetailResponse, error)
	UpdateTicket(ctx context.Context, ticketID int64, req TicketUpdateRequest) error
	DeleteTicket(ctx context.Context, ticketID int64) error
	AssignTicket(ctx context.Context, ticketID int64, req TicketAssignRequest) error
	SetPriority(ctx context.Context, ticketID int64, req TicketPriorityRequest) error
	CreateResolution(ctx context.Context, ticketID int64, req TicketResolutionRequest, file *storage.File) error
	CloseTicket(ctx context.Context, ticketID int64) error
}

type categorySvc interface {
	GetCategory(ctx context.Context, id int64) (*category.CategoryResponse, error)
}

type divisionSvc interface {
	GetDivision(ctx context.Context, id int64) (*division.DivisionResponse, error)
}

type service struct {
	repo            TicketRepository
	storage         storage.Storage
	storageConfig   config.Storage
	cache           cache.Cache
	notifier        mailer.Notifier
	notificationSvc notification.Notifier
	categorySvc     categorySvc
	divisionSvc     divisionSvc
}

func NewTicketService(
	repo TicketRepository,
	store storage.Storage,
	storageConfig config.Storage,
	cache cache.Cache,
	notifier mailer.Notifier,
	notificationSvc notification.Notifier,
	categorySvc categorySvc,
	divisionSvc divisionSvc,
) TicketService {
	return &service{
		repo:            repo,
		storage:         store,
		storageConfig:   storageConfig,
		cache:           cache,
		notifier:        notifier,
		notificationSvc: notificationSvc,
		categorySvc:     categorySvc,
		divisionSvc:     divisionSvc,
	}
}

func (s *service) ListTickets(ctx context.Context, params *GetTicketParams) ([]TicketResponse, int64, error) {
	if params == nil {
		params = &GetTicketParams{}
	}
	params.Normalize()

	tickets, total, err := s.repo.GetAll(ctx, *params)
	if err != nil {
		return nil, 0, err
	}

	return toTicketResponses(tickets), total, nil
}

func (s *service) CreateTicket(ctx context.Context, req *TicketCreateRequest, file *storage.File) error {
	if _, err := s.categorySvc.GetCategory(ctx, req.CategoryID); err != nil {
		return err
	}

	if _, err := s.divisionSvc.GetDivision(ctx, req.DivisionID); err != nil {
		return err
	}

	createdBy, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	ticket := Ticket{
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		DivisionID:  req.DivisionID,
		CreatedBy:   createdBy,
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				slog.ErrorContext(ctx, "rollback failed", "error", rbErr)
			}
		}
	}()

	ticketID, err := tx.Create(ctx, ticket)
	if err != nil {
		return err
	}

	if file != nil {
		objectKey := httputil.GenerateObjectKey(fmt.Sprintf("tickets/%d/report", ticketID), file.Filename)

		if uploadErr := s.storage.Upload(ctx, objectKey, file); uploadErr == nil {
			attachment := TicketAttachment{
				TicketID:       ticketID,
				FileKey:        objectKey,
				AttachmentType: string(Report),
				UploadedBy:     createdBy,
			}

			if attachErr := tx.CreateAttachment(ctx, attachment); attachErr != nil {
				_ = s.storage.Delete(ctx, objectKey)
				slog.ErrorContext(ctx, "failed to create ticket attachment", "ticket_id", ticketID, "error", attachErr)
			}
		} else {
			slog.ErrorContext(ctx, "failed to upload ticket attachment", "ticket_id", ticketID, "error", uploadErr)
		}
	}

	err = tx.Commit()
	if err == nil {
		dashboard.InvalidateCache(ctx, s.cache)
		s.notifier.NewTicketEmail(ctx, ticketID, req.Title, req.Description, createdBy)
		s.notificationSvc.NewTicket(ctx, ticketID, createdBy)
	}

	return err
}

func (s *service) GetTicket(ctx context.Context, id int64) (*TicketDetailResponse, error) {
	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return nil, apperr.NotFound("ticket")
		}
		return nil, err
	}

	attachments, err := s.repo.GetAttachmentsByTicketID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := toTicketDetailResponse(*ticket, attachments, s.storageConfig)

	return &result, nil
}

func (s *service) UpdateTicket(ctx context.Context, ticketID int64, req TicketUpdateRequest) error {
	existing, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperr.NotFound("ticket")
		}
		return err
	}

	if _, err := s.categorySvc.GetCategory(ctx, req.CategoryID); err != nil {
		return err
	}

	if _, err := s.divisionSvc.GetDivision(ctx, req.DivisionID); err != nil {
		return err
	}

	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	if existing.CreatedByID != userID {
		return apperr.Forbidden("you can only edit your own tickets")
	}

	if TicketStatus(existing.Status) != StatusOpen {
		return apperr.BadRequest("only open tickets can be edited")
	}

	if err := s.repo.Update(ctx, ticketID, Ticket{
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		DivisionID:  req.DivisionID,
		CreatedBy:   userID,
	}); err != nil {
		return err
	}

	dashboard.InvalidateCache(ctx, s.cache)

	return nil
}

func (s *service) DeleteTicket(ctx context.Context, ticketID int64) error {
	existing, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperr.NotFound("ticket")
		}
		return err
	}

	if TicketStatus(existing.Status) != StatusOpen {
		return apperr.BadRequest("only open tickets can be deleted")
	}

	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	if existing.CreatedByID != userID {
		return apperr.Forbidden("you can only delete your own tickets")
	}

	attachments, err := s.repo.GetAttachmentsByTicketID(ctx, ticketID)
	if err != nil {
		return err
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				slog.ErrorContext(ctx, "rollback failed", "error", rbErr)
			}
		}
	}()

	if err = tx.DeleteAttachmentsByTicketID(ctx, ticketID); err != nil {
		return err
	}

	if err = tx.Delete(ctx, ticketID); err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperr.NotFound("ticket")
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	dashboard.InvalidateCache(ctx, s.cache)

	if attachments != nil {
		for _, a := range *attachments {
			if delErr := s.storage.Delete(ctx, a.FileKey); delErr != nil {
				slog.ErrorContext(ctx, "failed to delete attachment from storage", "key", a.FileKey, "error", delErr)
			}
		}
	}

	return nil
}

func (s *service) AssignTicket(ctx context.Context, ticketID int64, req TicketAssignRequest) error {
	existing, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperr.NotFound("ticket")
		}

		return err
	}

	if existing.Priority == nil {
		return apperr.BadRequest("please set priority before assigning a ticket")
	}

	actorID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	if err := s.repo.Assign(ctx, ticketID, req.AssignedTo, actorID, req.Note); err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperr.NotFound("ticket")
		}

		if errors.Is(err, user.ErrUserNotFound) {
			return apperr.NotFound("user")
		}
		return err
	}

	dashboard.InvalidateCache(ctx, s.cache)
	s.notificationSvc.TicketAssigned(ctx, ticketID, req.AssignedTo, actorID)
	s.notificationSvc.TicketInProgress(ctx, ticketID, existing.CreatedByID, actorID)

	return nil
}

func (s *service) SetPriority(ctx context.Context, ticketID int64, req TicketPriorityRequest) error {
	if err := s.repo.UpdatePriority(ctx, ticketID, req.Priority); err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperr.NotFound("ticket")
		}

		return err
	}

	dashboard.InvalidateCache(ctx, s.cache)

	return nil
}

func (s *service) CreateResolution(ctx context.Context, ticketID int64, req TicketResolutionRequest, file *storage.File) error {
	existing, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperr.NotFound("ticket")
		}

		return err
	}

	if existing.AssignedToID == nil {
		return apperr.BadRequest("please assign the ticket before adding a resolution")
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				slog.ErrorContext(ctx, "rollback failed", "error", rbErr)
			}
		}
	}()

	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	if err = tx.UpdateResolution(ctx, ticketID, userID, req.Resolution); err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperr.NotFound("ticket")
		}
		return err
	}

	if file != nil {
		objectKey := httputil.GenerateObjectKey(fmt.Sprintf("tickets/%d/resolution", ticketID), file.Filename)

		if uploadErr := s.storage.Upload(ctx, objectKey, file); uploadErr == nil {
			attachment := TicketAttachment{
				TicketID:       ticketID,
				FileKey:        objectKey,
				AttachmentType: string(Resolution),
				UploadedBy:     userID,
			}

			if attachErr := tx.CreateAttachment(ctx, attachment); attachErr != nil {
				_ = s.storage.Delete(ctx, objectKey)
				slog.ErrorContext(ctx, "failed to create ticket attachment", "ticket_id", ticketID, "error", attachErr)
			}
		} else {
			slog.ErrorContext(ctx, "failed to upload ticket attachment", "ticket_id", ticketID, "error", uploadErr)
		}
	}

	err = tx.Commit()
	if err == nil {
		dashboard.InvalidateCache(ctx, s.cache)
	}

	return err
}

func (s *service) CloseTicket(ctx context.Context, ticketID int64) error {
	existing, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperr.NotFound("ticket")
		}
		return err
	}

	if TicketStatus(existing.Status) != StatusResolved {
		return apperr.BadRequest("only resolved tickets can be closed")
	}

	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	if existing.CreatedByID != userID {
		return apperr.Forbidden("you can only close your own tickets")
	}

	if err := s.repo.CloseTicket(ctx, ticketID, userID); err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperr.NotFound("ticket")
		}

		return err
	}

	dashboard.InvalidateCache(ctx, s.cache)
	s.notificationSvc.TicketClosed(ctx, ticketID, existing.CreatedByID, userID)

	return nil
}
