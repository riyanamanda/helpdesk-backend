package ticket

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/request"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/upload"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/validation"
)

type handler struct {
	svc TicketService
}

func NewTicketHandler(svc TicketService) *handler {
	return &handler{
		svc: svc,
	}
}

func (h *handler) ListTickets(c *echo.Context) error {
	var param GetTicketParams

	if err := c.Bind(&param); err != nil {
		return response.Error(c, apperror.BadRequest("invalid query params"))
	}

	tickets, total, err := h.svc.FetchAllTickets(c.Request().Context(), &param)
	if err != nil {
		return response.Error(c, err)
	}

	return response.WithPagination(c, http.StatusOK, tickets, param.Page, param.Limit, total)
}

func (h *handler) CreateTicket(c *echo.Context) error {
	req, err := request.BindAndValidate[TicketCreateRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	var file *upload.File

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
			return response.Error(c, apperror.Internal("failed to open uploaded attachment"))
		}
		defer func() {
			if err := f.Close(); err != nil {
				slog.Error("failed to close file", "error", err)
			}
		}()

		file = &upload.File{
			Content:     f,
			Filename:    fileHeader.Filename,
			ContentType: fileHeader.Header.Get("Content-Type"),
			Size:        fileHeader.Size,
		}
	}

	if err := h.svc.RegisterTicket(c.Request().Context(), req, file); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusCreated, "ticket created successfully")
}

func (h *handler) GetTicket(c *echo.Context) error {
	id, err := utils.ParsePositiveInt64PathParam(c, "id", "ticket")
	if err != nil {
		return response.Error(c, err)
	}

	ticket, err := h.svc.FindTicketByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, ticket)
}

func (h *handler) AssignTicket(c *echo.Context) error {
	ticketID, err := utils.ParsePositiveInt64PathParam(c, "id", "ticket")
	if err != nil {
		return response.Error(c, err)
	}

	req, err := request.BindAndValidate[TicketAssignRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.AssignTicket(c.Request().Context(), ticketID, *req); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusOK, "ticket assigned successfully")
}

func (h *handler) SetPriority(c *echo.Context) error {
	ticketID, err := utils.ParsePositiveInt64PathParam(c, "id", "ticket")
	if err != nil {
		return response.Error(c, err)
	}

	req, err := request.BindAndValidate[TicketPriorityRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.SetPriority(c.Request().Context(), ticketID, *req); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusOK, "priority has been set successfully")
}

func (h *handler) CreateResolution(c *echo.Context) error {
	ticketID, err := utils.ParsePositiveInt64PathParam(c, "id", "ticket")
	if err != nil {
		return response.Error(c, err)
	}

	req, err := request.BindAndValidate[TicketResolutionRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	var file *upload.File

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
			return response.Error(c, apperror.Internal("failed to open uploaded attachment"))
		}
		defer func() {
			if err := f.Close(); err != nil {
				slog.Error("failed to close file", "error", err)
			}
		}()

		file = &upload.File{
			Content:     f,
			Filename:    fileHeader.Filename,
			ContentType: fileHeader.Header.Get("Content-Type"),
			Size:        fileHeader.Size,
		}
	}

	if err := h.svc.RegisterResolution(c.Request().Context(), ticketID, *req, file); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusCreated, "resolution created successfully")
}

func (h *handler) CloseTicket(c *echo.Context) error {
	ticketID, err := utils.ParsePositiveInt64PathParam(c, "id", "ticket")
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.CloseTicket(c.Request().Context(), ticketID); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusOK, "ticket closed successfully")
}
