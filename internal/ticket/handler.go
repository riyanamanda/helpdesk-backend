package ticket

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/validation"
)

type Handler struct {
	svc TicketService
}

func NewTicketHandler(svc TicketService) *Handler {
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

	var file *storage.File

	fileHeader, err := c.FormFile("attachment")
	if err != nil && !errors.Is(err, http.ErrMissingFile) && !errors.Is(err, http.ErrNotMultipart) {
		return response.Error(c, err)
	}

	if fileHeader != nil {
		if err := validation.ValidateImage(fileHeader, maxTicketAttachmentSize, AllowedTicketAttachmentTypes); err != nil {
			return response.Error(c, err)
		}

		f, err := fileHeader.Open()
		if err != nil {
			return response.Error(c, apperr.Internal())
		}
		defer f.Close()

		file = &storage.File{
			Content:     f,
			Filename:    fileHeader.Filename,
			ContentType: fileHeader.Header.Get("Content-Type"),
			Size:        fileHeader.Size,
		}
	}

	if err := h.svc.CreateTicket(c.Request().Context(), req, file); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
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

	if err := h.svc.UpdateTicket(c.Request().Context(), ticketID, *req); err != nil {
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

	if err := h.svc.AssignTicket(c.Request().Context(), ticketID, *req); err != nil {
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

	if err := h.svc.SetPriority(c.Request().Context(), ticketID, *req); err != nil {
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

	var file *storage.File
	fileHeader, err := c.FormFile("attachment")
	if err != nil && !errors.Is(err, http.ErrMissingFile) && !errors.Is(err, http.ErrNotMultipart) {
		return response.Error(c, err)
	}

	if fileHeader != nil {
		if err := validation.ValidateImage(fileHeader, maxTicketAttachmentSize, AllowedTicketAttachmentTypes); err != nil {
			return response.Error(c, err)
		}

		f, err := fileHeader.Open()
		if err != nil {
			return response.Error(c, apperr.Internal())
		}
		defer f.Close()

		file = &storage.File{
			Content:     f,
			Filename:    fileHeader.Filename,
			ContentType: fileHeader.Header.Get("Content-Type"),
			Size:        fileHeader.Size,
		}
	}

	if err := h.svc.CreateResolution(c.Request().Context(), ticketID, *req, file); err != nil {
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
