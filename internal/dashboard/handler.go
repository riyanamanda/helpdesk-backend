package dashboard

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type handler struct {
	service DashboardService
}

func NewDashboardHandler(service DashboardService) handler {
	return handler{service: service}
}

func (h *handler) GetSummary(c *echo.Context) error {
	summary, err := h.service.GetSummary(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, summary)
}

func (h *handler) GetRecentTickets(c *echo.Context) error {
	tickets, err := h.service.GetRecentTickets(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, tickets)
}
