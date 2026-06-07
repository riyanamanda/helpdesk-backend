package notification

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type Handler struct {
	svc NotificationService
}

func NewNotificationHandler(svc NotificationService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) ListNotifications(c *echo.Context) error {
	notifications, err := h.svc.ListNotifications(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, notifications)
}

func (h *Handler) UnreadCount(c *echo.Context) error {
	resp, err := h.svc.UnreadCount(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, resp)
}

func (h *Handler) MarkAllAsRead(c *echo.Context) error {
	if err := h.svc.MarkAllAsRead(c.Request().Context()); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) MarkAsRead(c *echo.Context) error {
	id, err := httputil.ParsePositiveInt64PathParam(c, "id", "notification")
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.MarkAsRead(c.Request().Context(), id); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}
