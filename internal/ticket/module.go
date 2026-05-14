package ticket

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Group, db *sqlx.DB) {
	repo := NewTicketRepository(db)
	svc := NewTicketService(repo)
	handler := NewTicketHandler(svc)

	e.GET("/tickets", handler.ListTickets)
}
