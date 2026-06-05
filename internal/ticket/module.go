package ticket

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/notification"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/storage"
)

func Register(e *echo.Group, db *sqlx.DB, storageService storage.Storage, storageConfig config.Storage, cache cache.Cache, notifier notification.NotificationService) {
	repo := NewTicketRepository(db)
	svc := NewTicketService(repo, storageService, storageConfig, cache, notifier)
	handler := NewTicketHandler(svc)

	adminOnly := middleware.RequireRole("ADMIN")

	e.GET("/tickets", handler.ListTickets)
	e.POST("/tickets", handler.CreateTicket)
	e.GET("/tickets/:id", handler.GetTicket)
	e.PUT("/tickets/:id", handler.UpdateTicket)
	e.DELETE("/tickets/:id", handler.DeleteTicket)
	e.PATCH("/tickets/:id/assign", handler.AssignTicket, adminOnly)
	e.PATCH("/tickets/:id/priority", handler.SetPriority, adminOnly)
	e.PATCH("/tickets/:id/resolution", handler.CreateResolution, adminOnly)
	e.PATCH("/tickets/:id/close", handler.CloseTicket)
}
