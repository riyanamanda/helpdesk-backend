package dashboard

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
)

func Register(e *echo.Group, db *sqlx.DB, cache cache.Cache) {
	repo := NewDashboardRepository(db)
	svc := NewDashboardService(repo, cache)
	handler := NewDashboardHandler(svc)

	e.GET("/dashboard/summary", handler.GetSummary)
	e.GET("/dashboard/recent-tickets", handler.GetRecentTickets)
	e.GET("/dashboard/monthly-trend", handler.GetMonthlyTrend)
	e.GET("/dashboard/tickets-by-category", handler.GetTicketsByCategory)
	e.GET("/dashboard/agent-workload", handler.GetAgentWorkload)
}
