package ticket

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/validation"
)

type service interface {
	ListTickets(ctx context.Context, params *GetTicketParams) ([]TicketResponse, int64, error)
	CreateTicket(ctx context.Context, req *TicketCreateRequest, file *storage.File) error
	GetTicket(ctx context.Context, id int64) (*TicketDetailResponse, error)
	UpdateTicket(ctx context.Context, ticketID int64, req *TicketUpdateRequest) error // Disesuaikan ke Pointer
	DeleteTicket(ctx context.Context, ticketID int64) error
	AssignTicket(ctx context.Context, ticketID int64, req *TicketAssignRequest) error                             // Disesuaikan ke Pointer
	SetPriority(ctx context.Context, ticketID int64, req *TicketPriorityRequest) error                            // Disesuaikan ke Pointer
	CreateResolution(ctx context.Context, ticketID int64, req *TicketResolutionRequest, file *storage.File) error // Disesuaikan ke Pointer
	CloseTicket(ctx context.Context, ticketID int64) error
}

type Handler struct {
	svc service
}

func NewHandler(svc service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) ListTickets(c *echo.Context) error {
	var params GetTicketParams
	if err := c.Bind(&params); err != nil {
		return response.Error(c, apperr.BadRequest("invalid query params"))
	}

	tickets, total, err := h.svc.ListTickets(c.Request().Context(), &params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Paginated(c, tickets, params.Page, params.Limit, total)
}

func (h *Handler) CreateTicket(c *echo.Context) error {
	req, err := httputil.BindAndValidate[TicketCreateRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	file, cleanup, err := parseAttachment(c, "attachment")
	if err != nil {
		return response.Error(c, err)
	}
	if cleanup != nil {
		defer cleanup()
	}

	if err := h.svc.CreateTicket(c.Request().Context(), req, file); err != nil {
		return response.Error(c, err)
	}

	return response.Created(c)
}

func (h *Handler) GetTicket(c *echo.Context) error {
	id, err := httputil.ParsePositiveInt64PathParam(c, "id", "ticket")
	if err != nil {
		return response.Error(c, err)
	}

	ticket, err := h.svc.GetTicket(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, ticket)
}

func (h *Handler) UpdateTicket(c *echo.Context) error {
	ticketID, err := httputil.ParsePositiveInt64PathParam(c, "id", "ticket")
	if err != nil {
		return response.Error(c, err)
	}

	req, err := httputil.BindAndValidate[TicketUpdateRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.UpdateTicket(c.Request().Context(), ticketID, req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) DeleteTicket(c *echo.Context) error {
	ticketID, err := httputil.ParsePositiveInt64PathParam(c, "id", "ticket")
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.DeleteTicket(c.Request().Context(), ticketID); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) AssignTicket(c *echo.Context) error {
	ticketID, err := httputil.ParsePositiveInt64PathParam(c, "id", "ticket")
	if err != nil {
		return response.Error(c, err)
	}

	req, err := httputil.BindAndValidate[TicketAssignRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.AssignTicket(c.Request().Context(), ticketID, req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) SetPriority(c *echo.Context) error {
	ticketID, err := httputil.ParsePositiveInt64PathParam(c, "id", "ticket")
	if err != nil {
		return response.Error(c, err)
	}

	req, err := httputil.BindAndValidate[TicketPriorityRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.SetPriority(c.Request().Context(), ticketID, req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) CreateResolution(c *echo.Context) error {
	ticketID, err := httputil.ParsePositiveInt64PathParam(c, "id", "ticket")
	if err != nil {
		return response.Error(c, err)
	}

	req, err := httputil.BindAndValidate[TicketResolutionRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	file, cleanup, err := parseAttachment(c, "attachment")
	if err != nil {
		return response.Error(c, err)
	}
	if cleanup != nil {
		defer cleanup()
	}

	if err := h.svc.CreateResolution(c.Request().Context(), ticketID, req, file); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) CloseTicket(c *echo.Context) error {
	ticketID, err := httputil.ParsePositiveInt64PathParam(c, "id", "ticket")
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.CloseTicket(c.Request().Context(), ticketID); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

// Helper re-usable untuk ekstraksi file attachment dari Multipart Form
func parseAttachment(c *echo.Context, formName string) (*storage.File, func(), error) {
	fileHeader, err := c.FormFile(formName)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) || errors.Is(err, http.ErrNotMultipart) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	if err := validation.ValidateImage(fileHeader, maxTicketAttachmentSize, AllowedTicketAttachmentTypes); err != nil {
		return nil, nil, err
	}

	f, err := fileHeader.Open()
	if err != nil {
		return nil, nil, apperr.Internal()
	}

	file := &storage.File{
		Content:     f,
		Filename:    fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
		Size:        fileHeader.Size,
	}

	cleanup := func() {
		_ = f.Close()
	}

	return file, cleanup, nil
}
