package feedback

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/notification"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
)

func Register(e *echo.Group, db *sqlx.DB, notificationNotifier notification.Notifier) {
	repo := NewFeedbackRepository(db)
	svc := NewFeedbackService(repo, notificationNotifier)
	handler := NewFeedbackHandler(svc)

	admin := e.Group("/admin")
	admin.GET("/feedbacks", handler.ListAllFeedbacks, middleware.RequirePermission(rbac.PermissionFeedbackView))

	e.GET("/feedbacks", handler.ListFeedbacks, middleware.RequirePermission(rbac.PermissionFeedbackView))
	e.POST("/feedbacks", handler.CreateFeedback, middleware.RequirePermission(rbac.PermissionFeedbackCreate))
	e.GET("/feedbacks/:id", handler.GetFeedback, middleware.RequirePermission(rbac.PermissionFeedbackView))
	e.PATCH("/feedbacks/:id/status", handler.UpdateFeedbackStatus, middleware.RequirePermission(rbac.PermissionFeedbackUpdate))
}
