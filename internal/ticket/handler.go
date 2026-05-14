package ticket

import (
	"net/http"

	"github.com/labstack/echo/v5"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
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
