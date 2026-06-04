package feedback

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/middleware"
)

func Register(e *echo.Group, db *sqlx.DB) {
	repo := NewFeedbackRepository(db)
	svc := NewFeedbackService(repo)
	handler := NewFeedbackHandler(svc)

	adminOnly := middleware.RequireRole("ADMIN")

	admin := e.Group("/admin")
	admin.GET("/feedbacks", handler.ListAllFeedbacks, adminOnly)

	e.GET("/feedbacks", handler.ListFeedbacks)
	e.POST("/feedbacks", handler.CreateFeedback)
	e.GET("/feedbacks/:id", handler.GetFeedback)
	e.PATCH("/feedbacks/:id/status", handler.UpdateFeedbackStatus, adminOnly)
}
