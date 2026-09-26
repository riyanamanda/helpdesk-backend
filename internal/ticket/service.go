package ticket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/category"
	"github.com/riyanamanda/helpdesk-backend/internal/dashboard"
	"github.com/riyanamanda/helpdesk-backend/internal/division"
	"github.com/riyanamanda/helpdesk-backend/internal/event"
	"github.com/riyanamanda/helpdesk-backend/internal/outbox"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
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

type userSvc interface {
	GetUser(ctx context.Context, id uuid.UUID) (*user.UserResponse, error)
}

type service struct {
	repo          TicketRepository
	outboxRepo    outbox.Repository
	txManager     *database.Manager
	storage       storage.Storage
	storageConfig config.Storage
	cache         cache.Cache
	categorySvc   categorySvc
	divisionSvc   divisionSvc
	userSvc       userSvc
}

func NewTicketService(
	repo TicketRepository,
	outboxRepo outbox.Repository,
	txManager *database.Manager,
	store storage.Storage,
	storageConfig config.Storage,
	cache cache.Cache,
	categorySvc categorySvc,
	divisionSvc divisionSvc,
	userSvc userSvc,
) TicketService {
	return &service{
		repo:          repo,
		outboxRepo:    outboxRepo,
		txManager:     txManager,
		storage:       store,
		storageConfig: storageConfig,
		cache:         cache,
		categorySvc:   categorySvc,
		divisionSvc:   divisionSvc,
		userSvc:       userSvc,
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

	getUser, err := s.userSvc.GetUser(ctx, createdBy)
	if err != nil {
		return err
	}

	ticket := Ticket{
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		DivisionID:  req.DivisionID,
		CreatedBy:   createdBy,
	}

	// Begin transaction
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ticketID, err := s.repo.Create(ctx, tx, ticket)
	if err != nil {
		return err
	}

	event := event.TicketCreatedEvent{
		TicketID:    ticketID,
		SubmittedBy: getUser.Name,
		Title:       req.Title,
		Description: req.Description,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	outboxEvent := outbox.OutboxEvent{
		EventType:   "ticket.created",
		AggregateID: strconv.FormatInt(ticketID, 10),
		Payload:     payload,
	}

	if err := s.outboxRepo.Create(ctx, tx, outboxEvent); err != nil {
		return err
	}

	var objectKey string
	var uploaded bool

	defer func() {
		if uploaded {
			if err := s.storage.Delete(ctx, objectKey); err != nil {
				slog.ErrorContext(ctx, "failed to cleanup ticket attachment", "object_key", objectKey, "error", err)
			}
		}
	}()

	if file != nil {
		objectKey = httputil.GenerateObjectKey(fmt.Sprintf("tickets/%d/report", ticketID), file.Filename)

		if err := s.storage.Upload(ctx, objectKey, file); err != nil {
			return err
		}

		uploaded = true
		attachment := TicketAttachment{
			TicketID:       ticketID,
			FileKey:        objectKey,
			AttachmentType: string(Report),
			UploadedBy:     createdBy,
		}

		if err := s.repo.CreateAttachment(ctx, tx, attachment); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	uploaded = false
	dashboard.InvalidateCache(ctx, s.cache)

	return nil
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

	attachments, err := s.repo.GetAttachmentsByTicketID(ctx, ticketID)
	if err != nil {
		return err
	}

	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = s.repo.DeleteAttachmentsByTicketID(ctx, tx, ticketID); err != nil {
		return err
	}

	if err = s.repo.Delete(ctx, tx, ticketID); err != nil {
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

	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	userID, ok := ctxkey.GetUserIDFromContext(ctx)
	if !ok {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized")
	}

	if err = s.repo.UpdateResolution(ctx, tx, ticketID, req.ResolvedBy, req.Resolution); err != nil {
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

			if attachErr := s.repo.CreateAttachment(ctx, tx, attachment); attachErr != nil {
				_ = s.storage.Delete(ctx, objectKey)
				slog.ErrorContext(ctx, "failed to create ticket attachment", "ticket_id", ticketID, "error", attachErr)
			}
		} else {
			slog.ErrorContext(ctx, "failed to upload ticket attachment", "ticket_id", ticketID, "error", uploadErr)
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	dashboard.InvalidateCache(ctx, s.cache)

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

	if err := s.repo.CloseTicket(ctx, ticketID, userID); err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			return apperr.NotFound("ticket")
		}

		return err
	}

	dashboard.InvalidateCache(ctx, s.cache)

	return nil
}
