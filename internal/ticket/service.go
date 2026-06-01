package ticket

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/riyanamanda/helpdesk-backend/internal/dashboard"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperror"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/upload"
	"github.com/riyanamanda/helpdesk-backend/internal/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type TicketService interface {
	ListTickets(ctx context.Context, params *GetTicketParams) ([]TicketResponse, int64, error)
	CreateTicket(ctx context.Context, req *TicketCreateRequest, file *upload.File) error
	GetTicket(ctx context.Context, id int64) (TicketDetailResponse, error)
	AssignTicket(ctx context.Context, ticketID int64, req TicketAssignRequest) error
	SetPriority(ctx context.Context, ticketID int64, req TicketPriorityRequest) error
	CreateResolution(ctx context.Context, ticketID int64, req TicketResolutionRequest, file *upload.File) error
	CloseTicket(ctx context.Context, ticketID int64) error
}

type service struct {
	repo          TicketRepository
	storage       storage.Storage
	storageConfig config.Storage
	cache         cache.Cache
}

func NewTicketService(repo TicketRepository, store storage.Storage, storageConfig config.Storage, cache cache.Cache) TicketService {
	return &service{
		repo:          repo,
		storage:       store,
		storageConfig: storageConfig,
		cache:         cache,
	}
}

func (s *service) ListTickets(ctx context.Context, params *GetTicketParams) ([]TicketResponse, int64, error) {
	if params == nil {
		params = &GetTicketParams{}
	}
	params.Normalize()

	tickets, total, err := s.repo.GetAll(ctx, *params)
	if err != nil {
		return []TicketResponse{}, 0, err
	}

	return toTicketResponses(tickets), total, nil
}

func (s *service) CreateTicket(ctx context.Context, req *TicketCreateRequest, file *upload.File) error {
	createdBy, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperror.Unauthorized(apperror.CodeUnauthorized, "unauthorized")
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

		if uploadErr := s.storage.Upload(ctx, objectKey, file.Content, file.Size, file.ContentType); uploadErr == nil {
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
	}

	return err
}

func (s *service) GetTicket(ctx context.Context, id int64) (TicketDetailResponse, error) {
	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return TicketDetailResponse{}, apperror.NotFound("ticket")
		}
		return TicketDetailResponse{}, err
	}

	attachments, err := s.repo.GetAttachmentsByTicketID(ctx, id)
	if err != nil {
		return TicketDetailResponse{}, err
	}

	return toTicketDetailResponse(*ticket, attachments, s.storageConfig), nil
}

func (s *service) AssignTicket(ctx context.Context, ticketID int64, req TicketAssignRequest) error {
	existing, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperror.NotFound("ticket")
		}

		return err
	}

	if existing.Priority == nil {
		return apperror.BadRequest("Please set priority before assign a ticket")
	}

	if err := s.repo.Assign(ctx, ticketID, req.AssignedTo); err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperror.NotFound("ticket")
		}

		if errors.Is(err, user.ErrUserNotFound) {
			return apperror.NotFound("user")
		}
		return err
	}

	dashboard.InvalidateCache(ctx, s.cache)

	return nil
}

func (s *service) SetPriority(ctx context.Context, ticketID int64, req TicketPriorityRequest) error {
	if err := s.repo.UpdatePriority(ctx, ticketID, req.Priority); err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperror.NotFound("ticket")
		}

		return err
	}

	dashboard.InvalidateCache(ctx, s.cache)

	return nil
}

func (s *service) CreateResolution(ctx context.Context, ticketID int64, req TicketResolutionRequest, file *upload.File) error {
	existing, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperror.NotFound("ticket")
		}

		return err
	}

	if existing.AssignedToID == nil {
		return apperror.BadRequest("Please assign the ticket before add resolution")
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
		return apperror.Unauthorized(apperror.CodeUnauthorized, "unauthorized")
	}

	if err = tx.UpdateResolution(ctx, ticketID, userID, req.Resolution); err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperror.NotFound("ticket")
		}
		return err
	}

	if file != nil {
		objectKey := httputil.GenerateObjectKey(fmt.Sprintf("tickets/%d/resolution", ticketID), file.Filename)

		if uploadErr := s.storage.Upload(ctx, objectKey, file.Content, file.Size, file.ContentType); uploadErr == nil {
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
	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperror.Unauthorized(apperror.CodeUnauthorized, "unauthorized")
	}

	if err := s.repo.CloseTicket(ctx, ticketID, userID); err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperror.NotFound("ticket")
		}

		return err
	}

	dashboard.InvalidateCache(ctx, s.cache)

	return nil
}
