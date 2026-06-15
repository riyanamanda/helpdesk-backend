package ticket

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
	"github.com/riyanamanda/helpdesk-backend/internal/notification"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func Register(e *echo.Group, db *sqlx.DB, storageService storage.Storage, storageConfig config.Storage, cache cache.Cache, notifier mailer.Notifier, userRepo user.UserRepository, notificationNotifier notification.Notifier) {
	repo := NewTicketRepository(db)
	svc := NewTicketService(repo, storageService, storageConfig, cache, notifier, notificationNotifier)
	handler := NewTicketHandler(svc)

	e.GET("/tickets", handler.ListTickets)
	e.POST("/tickets", handler.CreateTicket)
	e.GET("/tickets/:id", handler.GetTicket)
	e.PUT("/tickets/:id", handler.UpdateTicket)
	e.DELETE("/tickets/:id", handler.DeleteTicket)
	e.PATCH("/tickets/:id/assign", handler.AssignTicket, middleware.RequiredPermission("ticket:assign"))
	e.PATCH("/tickets/:id/priority", handler.SetPriority, middleware.RequiredPermission("ticket:priority"))
	e.PATCH("/tickets/:id/resolution", handler.CreateResolution, middleware.RequiredPermission("ticket:resolution"))
	e.PATCH("/tickets/:id/close", handler.CloseTicket)
}
