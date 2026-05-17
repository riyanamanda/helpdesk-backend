package ticket

import (
	"errors"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo/v5"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/request"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/validation"
)

type handler struct {
	service TicketService
}

func NewTicketHandler(service TicketService) handler {
	return handler{
		service: service,
	}
}

func (h *handler) ListTickets(c *echo.Context) error {
	var param GetTicketParams

	if err := c.Bind(&param); err != nil {
		return response.Error(c, apperror.BadRequest("invalid query params"))
	}

	tickets, total, err := h.service.FetchAllTickets(c.Request().Context(), &param)
	if err != nil {
		return response.Error(c, err)
	}

	return response.WithPagination(c, http.StatusOK, tickets, param.Page, param.Limit, total)
}

func (h *handler) CreateTicket(c *echo.Context) error {
	var (
		file       multipart.File
		fileHeader *multipart.FileHeader
	)

	req, err := request.BindAndValidate[TicketCreateRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	fileHeader, err = c.FormFile("attachment")
	if err != nil && !errors.Is(err, http.ErrMissingFile) && !errors.Is(err, http.ErrNotMultipart) {
		return response.Error(c, err)
	}

	if fileHeader != nil {
		if err := validation.ValidateImage(fileHeader, maxTicketAttachmentSize, AllowedTicketAttachmentTypes); err != nil {
			return response.Error(c, err)
		}

		file, err = fileHeader.Open()
		if err != nil {
			return response.Error(c, apperror.Internal("failed to open uploaded attachment"))
		}
		defer func() {
			if err := file.Close(); err != nil {
				slog.Error("failed to close file", "error", err)
			}
		}()
	}

	if err := h.service.RegisterTicket(c.Request().Context(), req, file, fileHeader); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusCreated, "ticket created successfully")
}

func (h *handler) GetTicket(c *echo.Context) error {
	id, err := utils.ParsePositiveInt64PathParam(c, "id", "ticket")
	if err != nil {
		return response.Error(c, err)
	}

	ticket, err := h.service.FindTicketByID(c.Request().Context(), id)
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

	if err := h.service.AssignTicket(c.Request().Context(), ticketID, *req); err != nil {
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

	if err := h.service.SetPriority(c.Request().Context(), ticketID, *req); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusOK, "priority has been set successfully")
}

func (h *handler) CreateResolution(c *echo.Context) error {
	var (
		file       multipart.File
		fileHeader *multipart.FileHeader
	)

	ticketID, err := utils.ParsePositiveInt64PathParam(c, "id", "ticket")
	if err != nil {
		return response.Error(c, err)
	}

	req, err := request.BindAndValidate[TicketResolutionRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	fileHeader, err = c.FormFile("attachment")
	if err != nil && !errors.Is(err, http.ErrMissingFile) && !errors.Is(err, http.ErrNotMultipart) {
		return response.Error(c, err)
	}

	if fileHeader != nil {
		if err := validation.ValidateImage(fileHeader, maxTicketAttachmentSize, AllowedTicketAttachmentTypes); err != nil {
			return response.Error(c, err)
		}

		file, err = fileHeader.Open()
		if err != nil {
			return response.Error(c, apperror.Internal("failed to open uploaded attachment"))
		}
		defer func() {
			if err := file.Close(); err != nil {
				slog.Error("failed to close file", "error", err)
			}
		}()
	}

	if err := h.service.RegisterResolution(c.Request().Context(), ticketID, *req, file, fileHeader); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusCreated, "resolution created successfully")
}
