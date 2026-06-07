package notification

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Group, db *sqlx.DB) {
	repo := NewNotificationRepository(db)
	svc := NewNotificationService(repo)
	handler := NewNotificationHandler(svc)

	e.GET("/notifications", handler.ListNotifications)
	e.GET("/notifications/unread-count", handler.UnreadCount)
	e.PATCH("/notifications/read-all", handler.MarkAllAsRead)
	e.PATCH("/notifications/:id/read", handler.MarkAsRead)
}
