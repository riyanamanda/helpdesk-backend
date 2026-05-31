package ticket

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/storage"
)

func Register(e *echo.Group, db *sqlx.DB, store storage.Storage, storageConfig config.Storage, cache cache.Cache) {
	repo := NewTicketRepository(db)
	svc := NewTicketService(repo, store, storageConfig, cache)
	handler := NewTicketHandler(svc)

	e.GET("/tickets", handler.ListTickets)
	e.POST("/tickets", handler.CreateTicket)
	e.GET("/tickets/:id", handler.GetTicket)
	e.PATCH("/tickets/:id/assign", handler.AssignTicket)
	e.PATCH("/tickets/:id/priority", handler.SetPriority)
	e.PATCH("/tickets/:id/resolution", handler.CreateResolution)
	e.PATCH("/tickets/:id/close", handler.CloseTicket)
}
