package ticket

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/category"
	"github.com/riyanamanda/helpdesk-backend/internal/division"
	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
	"github.com/riyanamanda/helpdesk-backend/internal/notification"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func Register(e *echo.Group, db *sqlx.DB, storageService storage.Storage, storageConfig config.Storage, cache cache.Cache, notifier mailer.Notifier, userRepo user.UserRepository, notificationNotifier notification.Notifier) {
	repo := NewTicketRepository(db)
	catRepo := category.NewCategoryRepository(db)
	divRepo := division.NewDivisionRepository(db)

	svc := NewTicketService(repo, storageService, storageConfig, cache, notifier, notificationNotifier, catRepo, divRepo)
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
