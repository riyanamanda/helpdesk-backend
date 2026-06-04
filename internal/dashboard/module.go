package dashboard

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
)

func Register(e *echo.Group, db *sqlx.DB, cache cache.Cache) {
	repo := NewDashboardRepository(db)
	svc := NewDashboardService(repo, cache)
	handler := NewDashboardHandler(svc)

	adminOnly := middleware.RequireRole("ADMIN")

	e.GET("/dashboard/summary", handler.GetSummary, adminOnly)
	e.GET("/dashboard/recent-tickets", handler.GetRecentTickets, adminOnly)
}
