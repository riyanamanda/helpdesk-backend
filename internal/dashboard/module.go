package dashboard

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Group, db *sqlx.DB) {
	repo := NewDashboardRepository(db)
	svc := NewDashboardService(repo)
	handler := NewDashboardHandler(svc)

	e.GET("/dashboard/summary", handler.GetSummary)
	e.GET("/dashboard/recent-tickets", handler.GetRecentTickets)
}
