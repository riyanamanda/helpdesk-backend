package dashboard

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Group, db *sqlx.DB) {
	repo := NewDashboardRepository(db)
	svc := NewDashboardService(repo)
	h := NewDashboardHandler(svc)

	e.GET("/dashboard/summary", h.GetSummary)
	e.GET("/dashboard/recent-tickets", h.GetRecentTickets)
}
