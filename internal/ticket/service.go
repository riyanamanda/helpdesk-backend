package ticket

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/upload"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
	"github.com/riyanamanda/helpdesk-backend/internal/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type TicketService interface {
	FetchAllTickets(ctx context.Context, params *GetTicketParams) ([]TicketResponse, int64, error)
	RegisterTicket(ctx context.Context, req *TicketCreateRequest, file *upload.File) error
	FindTicketByID(ctx context.Context, id int64) (TicketDetailResponse, error)
	AssignTicket(ctx context.Context, ticketID int64, req TicketAssignRequest) error
	SetPriority(ctx context.Context, ticketID int64, req TicketPriorityRequest) error
	RegisterResolution(ctx context.Context, ticketID int64, req TicketResolutionRequest, file *upload.File) error
	CloseTicket(ctx context.Context, ticketID int64) error
}

type service struct {
	repo          TicketRepository
	storage       storage.Storage
	storageConfig config.Storage
}

func NewTicketService(repo TicketRepository, storage storage.Storage, storageConfig config.Storage) TicketService {
	return &service{
		repo:          repo,
		storage:       storage,
		storageConfig: storageConfig,
	}
}

func (s *service) FetchAllTickets(ctx context.Context, params *GetTicketParams) ([]TicketResponse, int64, error) {
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

func (s *service) RegisterTicket(ctx context.Context, req *TicketCreateRequest, file *upload.File) error {
	var createdBy uuid.UUID
	if currentUserID, ok := utils.GetUserIDFromContext(ctx); ok {
		createdBy = currentUserID
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
			_ = tx.Rollback()
		}
	}()

	ticketID, err := tx.Create(ctx, ticket)
	if err != nil {
		return err
	}

	if file != nil {
		objectKey := utils.GenerateObjectKey(fmt.Sprintf("tickets/%d/report", ticketID), file.Filename)

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
	return err
}

func (s *service) FindTicketByID(ctx context.Context, id int64) (TicketDetailResponse, error) {
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
	if err := s.repo.Assign(ctx, ticketID, req.AssignedTo); err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperror.NotFound("ticket")
		}
		if errors.Is(err, user.ErrUserNotFound) {
			return apperror.NotFound("user")
		}
		return err
	}

	return nil
}

func (s *service) SetPriority(ctx context.Context, ticketID int64, req TicketPriorityRequest) error {
	if err := s.repo.UpdatePriority(ctx, ticketID, req.Priority); err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperror.NotFound("ticket")
		}
		return err
	}

	return nil
}

func (s *service) RegisterResolution(ctx context.Context, ticketID int64, req TicketResolutionRequest, file *upload.File) error {
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var userID uuid.UUID
	if currentUserID, ok := utils.GetUserIDFromContext(ctx); ok {
		userID = currentUserID
	}

	if err = tx.UpdateResolution(ctx, ticketID, userID, req.Resolution); err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperror.NotFound("ticket")
		}
		return err
	}

	if file != nil {
		objectKey := utils.GenerateObjectKey(fmt.Sprintf("tickets/%d/resolution", ticketID), file.Filename)

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
	return err
}

func (s *service) CloseTicket(ctx context.Context, ticketID int64) error {
	var userID uuid.UUID
	if currentUser, ok := utils.GetUserIDFromContext(ctx); ok {
		userID = currentUser
	}

	if err := s.repo.CloseTicket(ctx, ticketID, userID); err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperror.NotFound("ticket")
		}
		return err
	}

	return nil
}
