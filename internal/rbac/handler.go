package rbac

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type service interface {
	ListRoles(ctx context.Context) ([]RoleResponse, error)
	ListPermissions(ctx context.Context) ([]PermissionResponse, error)
	GetRolePermissions(ctx context.Context, roleID int64) ([]PermissionResponse, error)
	SetRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error
}

type Handler struct {
	svc service
}

func NewHandler(svc service) *Handler {
	return &Handler{
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
	permissions, err := h.svc.ListPermissions(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, permissions)
}

func (h *Handler) GetRolePermissions(c *echo.Context) error {
	roleID, err := httputil.ParsePositiveInt64PathParam(c, "id", "role")
	if err != nil {
		return response.Error(c, apperr.BadRequest("invalid role id"))
	}

	permissions, err := h.svc.GetRolePermissions(c.Request().Context(), roleID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, permissions)
}

func (h *Handler) SetRolePermissions(c *echo.Context) error {
	roleID, err := httputil.ParsePositiveInt64PathParam(c, "id", "role")
	if err != nil {
		return response.Error(c, apperr.BadRequest("invalid role id"))
	}

	req, err := httputil.BindAndValidate[SetRolePermissionsRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.SetRolePermissions(c.Request().Context(), roleID, req.PermissionIDs); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}
