package rbac

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type Handler struct {
	svc RBACService
}

func NewRBACHandler(svc RBACService) Handler {
	return Handler{
		svc: svc,
	}
}

func (h *Handler) ListRoles(c *echo.Context) error {
	roles, err := h.svc.ListRoles(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, roles)
}

func (h *Handler) ListPermissions(c *echo.Context) error {
	roles, err := h.svc.ListPermissions(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, roles)
}
