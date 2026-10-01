package dashboard

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type service interface {
	GetSummary(ctx context.Context) (*SummaryResponse, error)
	GetMonthlyTrend(ctx context.Context, year int) ([]MonthlyTrendResponse, error)
	GetTicketsByCategory(ctx context.Context) ([]CategoryTicketsResponse, error)
	GetAgentWorkload(ctx context.Context) ([]AgentWorkloadResponse, error)
}

type Handler struct {
	svc service
}

func NewDashboardHandler(svc service) *Handler {
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

func (h *Handler) GetTicketsByCategory(c *echo.Context) error {
	categories, err := h.svc.GetTicketsByCategory(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, categories)
}

func (h *Handler) GetAgentWorkload(c *echo.Context) error {
	workload, err := h.svc.GetAgentWorkload(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, workload)
}
