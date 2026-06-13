package user

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type Handler struct {
	svc UserService
}

func NewUserHandler(svc UserService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) ListUsers(c *echo.Context) error {
	var params GetUserParams

	if err := c.Bind(&params); err != nil {
		return response.Error(c, apperr.BadRequest("invalid query params"))
	}

	users, total, err := h.svc.ListUsers(c.Request().Context(), &params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Paginated(c, users, params.Page, params.Limit, total)
}

func (h *Handler) ListAssignableUser(c *echo.Context) error {
	users, err := h.svc.ListAssignableUser(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, users)
}

func (h *Handler) CreateUser(c *echo.Context) error {
	req, err := httputil.BindAndValidate[UserCreateRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.CreateUser(c.Request().Context(), req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) GetUser(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, apperr.BadRequest("invalid user id"))
	}

	user, err := h.svc.GetUser(c.Request().Context(), &id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, user)
}

func (h *Handler) UpdateUser(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, apperr.BadRequest("invalid user id"))
	}

	req, err := httputil.BindAndValidate[UserUpdateRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	err = h.svc.UpdateUser(c.Request().Context(), id, req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) UpdatePassword(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, apperr.BadRequest("invalid user id"))
	}

	req, err := httputil.BindAndValidate[UserUpdatePasswordRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	err = h.svc.UpdatePassword(c.Request().Context(), id, req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}
