package ticket

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/storage"
)

func Register(e *echo.Group, db *sqlx.DB, storageService storage.Storage) {
	repo := NewTicketRepository(db)
	svc := NewTicketService(db, repo, storageService)
	handler := NewTicketHandler(svc)

	e.GET("/tickets", handler.ListTickets)
	e.POST("/tickets", handler.CreateTicket)
	e.GET("/tickets/:id", handler.GetTicket)
	e.PATCH("/tickets/:id/assign", handler.AssignTicket)
	e.PATCH("/tickets/:id/priority", handler.SetPriority)
	e.PATCH("/tickets/:id/resolution", handler.CreateResolution)
}
