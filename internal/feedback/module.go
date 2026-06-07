package feedback

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/notification"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func Register(e *echo.Group, db *sqlx.DB, userRepo user.UserRepository) {
	repo := NewFeedbackRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	notificationNotifier := notification.NewNotifier(notificationRepo, userRepo)
	svc := NewFeedbackService(repo, notificationNotifier)
	handler := NewFeedbackHandler(svc)

	adminOnly := middleware.RequireRole("ADMIN")

	admin := e.Group("/admin")
	admin.GET("/feedbacks", handler.ListAllFeedbacks, adminOnly)

	e.GET("/feedbacks", handler.ListFeedbacks)
	e.POST("/feedbacks", handler.CreateFeedback)
	e.GET("/feedbacks/:id", handler.GetFeedback)
	e.PATCH("/feedbacks/:id/status", handler.UpdateFeedbackStatus, adminOnly)
}
