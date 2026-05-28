package dashboard

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type handler struct {
	svc DashboardService
}

func NewDashboardHandler(svc DashboardService) *handler {
	return &handler{
		svc: svc,
	}
}

func (h *handler) GetSummary(c *echo.Context) error {
	summary, err := h.svc.GetSummary(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, summary)
}

func (h *handler) GetRecentTickets(c *echo.Context) error {
	tickets, err := h.svc.GetRecentTickets(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, tickets)
}
