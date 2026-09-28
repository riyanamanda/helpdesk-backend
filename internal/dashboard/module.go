package dashboard

import (
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Group, svc DashboardService) {
	handler := NewDashboardHandler(svc)

	e.GET("/dashboard/summary", handler.GetSummary)
	e.GET("/dashboard/monthly-trend", handler.GetMonthlyTrend)
	e.GET("/dashboard/tickets-by-category", handler.GetTicketsByCategory)
	e.GET("/dashboard/agent-workload", handler.GetAgentWorkload)
}
