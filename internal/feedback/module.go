package feedback

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/notification"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
)

func Register(e *echo.Group, db *sqlx.DB, notificationNotifier notification.Notifier) {
	repo := NewFeedbackRepository(db)
	svc := NewFeedbackService(repo, notificationNotifier)
	handler := NewFeedbackHandler(svc)

	admin := e.Group("/admin")
	admin.GET("/feedbacks", handler.ListAllFeedbacks, middleware.RequiredPermission("feedback:view"))

	e.GET("/feedbacks", handler.ListFeedbacks)
	e.POST("/feedbacks", handler.CreateFeedback)
	e.GET("/feedbacks/:id", handler.GetFeedback)
	e.PATCH("/feedbacks/:id/status", handler.UpdateFeedbackStatus, middleware.RequiredPermission("feedback:update"))
}
