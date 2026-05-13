package user

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/request"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type handler struct {
	svc UserService
}

func NewUserHandler(svc UserService) handler {
	return handler{
		svc: svc,
	}
}

func (h *handler) ListUser(c *echo.Context) error {
	var params GetUserParams

	if err := c.Bind(&params); err != nil {
		return response.Error(c, apperror.BadRequest("invalid query params"))
	}

	users, total, err := h.svc.GetUser(c.Request().Context(), &params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.WithPagination(c, http.StatusOK, users, params.Page, params.Limit, total)
}

func (h *handler) Create(c *echo.Context) error {
	req, err := request.BindAndValidate[UserCreateRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	user, err := h.svc.Create(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusCreated, user)
}

func (h *handler) GetByID(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, apperror.BadRequest("invalid user id"))
	}

	user, err := h.svc.GetById(c.Request().Context(), &id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, user)
}
