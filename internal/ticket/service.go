package ticket

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
	"github.com/riyanamanda/helpdesk-backend/internal/storage"
)

type TicketService interface {
	FetchAllTickets(ctx context.Context, params *GetTicketParams) ([]TicketResponse, int64, error)
	RegisterTicket(ctx context.Context, req *TicketCreateRequest, file multipart.File, fileHeader *multipart.FileHeader) error
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
		objectKey := fmt.Sprintf("tickets/%d/report/%d-%s", ticketID, time.Now().Unix(), fileHeader.Filename)
		contentType := fileHeader.Header.Get("Content-Type")

		err := s.storage.Upload(ctx, objectKey, file, fileHeader.Size, contentType)
		if err != nil {
			return err
		}

		attachment := TicketAttachment{
			TicketID:       ticketID,
			FileKey:        objectKey,
			AttachmentType: string(REPORT),
			UploadedBy:     createdBy,
		}

		if err := s.repo.CreateAttachment(ctx, attachment); err != nil {
			_ = s.storage.Delete(ctx, objectKey)
			return err
		}

	}

	return nil
}
