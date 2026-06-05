package dashboard

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type Handler struct {
	svc DashboardService
}

func NewDashboardHandler(svc DashboardService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) GetSummary(c *echo.Context) error {
	summary, err := h.svc.GetSummary(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, summary)
}

func (h *Handler) GetRecentTickets(c *echo.Context) error {
	tickets, err := h.svc.GetRecentTickets(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, tickets)
}

func (h *Handler) GetMonthlyTrend(c *echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 || year > time.Now().Year() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid year"})
	}

	trend, err := h.svc.GetMonthlyTrend(c.Request().Context(), year)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, trend)
}
