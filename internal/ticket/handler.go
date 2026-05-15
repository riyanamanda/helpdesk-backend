package ticket

import (
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo/v5"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/request"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
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
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
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
		defer file.Close()
	}

	if err := h.service.RegisterTicket(c.Request().Context(), req, file, fileHeader); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusCreated, "ticket created successfully")
}
