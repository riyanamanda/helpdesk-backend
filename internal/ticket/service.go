package ticket

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"

	"github.com/google/uuid"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
	"github.com/riyanamanda/helpdesk-backend/internal/storage"
)

type TicketService interface {
	FetchAllTickets(ctx context.Context, params *GetTicketParams) ([]TicketResponse, int64, error)
	RegisterTicket(ctx context.Context, req *TicketCreateRequest, file multipart.File, fileHeader *multipart.FileHeader) error
	FindTicketByID(ctx context.Context, id int64) (TicketDetailResponse, error)
	AssignTicket(ctx context.Context, ticketID int64, req TicketAssignRequest) error

	RegisterResolution(ctx context.Context, ticketID int64, req TicketResolutionCreateRequest, file multipart.File, fileHeader *multipart.FileHeader) error
	FindResolutionByTicketID(ctx context.Context, ticketID int64) (*TicketResolutionResponse, error)
}

type service struct {
	repo    TicketRepository
	storage storage.Storage
}

func NewTicketService(repo TicketRepository, storage storage.Storage) TicketService {
	return &service{
		repo:    repo,
		storage: storage,
	}
}

func (s *service) FetchAllTickets(ctx context.Context, params *GetTicketParams) ([]TicketResponse, int64, error) {
	if params == nil {
		params = &GetTicketParams{}
	}

	page, limit, _ := params.Normalize()
	params.Page = page
	params.Limit = limit

	tickets, total, err := s.repo.GetAll(ctx, *params)
	if err != nil {
		return []TicketResponse{}, 0, err
	}

	return toTicketResponses(tickets), total, nil
}

func (s *service) RegisterTicket(ctx context.Context, req *TicketCreateRequest, file multipart.File, fileHeader *multipart.FileHeader) error {
	var createdBy uuid.UUID
	if currentUserID, ok := utils.GetUserIDFromContext(ctx); ok {
		createdBy = currentUserID
	}

	ticket := Ticket{
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		CreatedBy:   createdBy,
	}

	ticketID, err := s.repo.Create(ctx, ticket)
	if err != nil {
		return err
	}

	if file != nil && fileHeader != nil {
		objectKey := utils.GenerateObjectKey(fmt.Sprintf("tickets/%d/report", ticketID), fileHeader.Filename)
		contentType := fileHeader.Header.Get("Content-Type")

		err := s.storage.Upload(ctx, objectKey, file, fileHeader.Size, contentType)

		// keep success response even though attachment error but remove uploaded file
		if err == nil {
			attachment := TicketAttachment{
				TicketID:       ticketID,
				FileKey:        objectKey,
				AttachmentType: string(REPORT),
				UploadedBy:     createdBy,
			}

			if err := s.repo.CreateAttachment(ctx, attachment); err != nil {

				_ = s.storage.Delete(ctx, objectKey)

				slog.ErrorContext(
					ctx,
					"failed to create ticket attachment",
					"ticket_id", ticketID,
					"error", err,
				)
			}
		} else {
			slog.ErrorContext(ctx, "failed to upload ticket attachment", "ticket_id", ticketID, "error", err)
		}

	}

	return nil
}

func (s *service) FindTicketByID(ctx context.Context, id int64) (TicketDetailResponse, error) {
	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return TicketDetailResponse{}, apperror.NotFound("ticket")
		}
		return TicketDetailResponse{}, err
	}

	attachment, err := s.repo.GetAttachmentByTicketID(ctx, id, REPORT)
	if err != nil {
		return TicketDetailResponse{}, err
	}

	return toTicketDetailResponse(*ticket, attachment, s.storage), nil
}

func (s *service) AssignTicket(ctx context.Context, ticketID int64, req TicketAssignRequest) error {

	if err := s.repo.Assign(ctx, ticketID, req.AssignedTo); err != nil {
		return err
	}

	return nil
}

func (s *service) RegisterResolution(ctx context.Context, ticketID int64, req TicketResolutionCreateRequest, file multipart.File, fileHeader *multipart.FileHeader) error {
	var userID uuid.UUID
	if currentUser, ok := utils.GetUserIDFromContext(ctx); ok {
		userID = currentUser
	}

	ticketResolution := TicketResolution{
		TicketID:   ticketID,
		ResolvedBy: userID,
		Resolution: req.Resolution,
	}

	if err := s.repo.CreateResolution(ctx, ticketResolution); err != nil {
		if errors.Is(err, ErrTicketResolutionAlreadyExists) {
			return apperror.AlreadyExists("ticket resolution")
		}
		return err
	}

	if file != nil && fileHeader != nil {
		objectKey := utils.GenerateObjectKey(fmt.Sprintf("tickets/%d/resolution", ticketID), fileHeader.Filename)
		contentType := fileHeader.Header.Get("Content-Type")

		err := s.storage.Upload(ctx, objectKey, file, fileHeader.Size, contentType)
		if err == nil {
			attachment := TicketAttachment{
				TicketID:       ticketID,
				FileKey:        objectKey,
				AttachmentType: string(RESOLUTION),
				UploadedBy:     userID,
			}

			if err := s.repo.CreateAttachment(ctx, attachment); err != nil {
				_ = s.storage.Delete(ctx, objectKey)

				slog.ErrorContext(ctx, "failed to create resolution attachment", "ticket_id", ticketID, "error", err)
			}
		} else {
			slog.ErrorContext(ctx, "failed to create resolution attachment", "ticket_id", ticketID, "error", err)
		}
	}

	return nil
}

func (s *service) FindResolutionByTicketID(ctx context.Context, ticketID int64) (*TicketResolutionResponse, error) {
	resolution, err := s.repo.GetResolutionByTicketID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	if resolution == nil {
		return nil, nil
	}

	response := toTicketResolutionResponse(*resolution)
	return &response, nil
}
