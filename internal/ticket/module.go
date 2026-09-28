package ticket

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
)

func Register(e *echo.Group, svc TicketService) {
	handler := NewTicketHandler(svc)

	e.GET("/tickets", handler.ListTickets, middleware.RequirePermission(rbac.PermissionTicketView))
	e.POST("/tickets", handler.CreateTicket, middleware.RequirePermission(rbac.PermissionTicketCreate))
	e.GET("/tickets/:id", handler.GetTicket, middleware.RequirePermission(rbac.PermissionTicketView))
	e.PUT("/tickets/:id", handler.UpdateTicket, middleware.RequirePermission(rbac.PermissionTicketUpdate))
	e.DELETE("/tickets/:id", handler.DeleteTicket, middleware.RequirePermission(rbac.PermissionTicketDelete))
	e.PATCH("/tickets/:id/assign", handler.AssignTicket, middleware.RequirePermission(rbac.PermissionTicketAssign))
	e.PATCH("/tickets/:id/priority", handler.SetPriority, middleware.RequirePermission(rbac.PermissionTicketPriority))
	e.PATCH("/tickets/:id/resolution", handler.CreateResolution, middleware.RequirePermission(rbac.PermissionTicketResolution))
	e.PATCH("/tickets/:id/close", handler.CloseTicket, middleware.RequirePermission(rbac.PermissionTicketClose))
}
