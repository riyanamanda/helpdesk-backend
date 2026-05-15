package user

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/request"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/validation"
)

type handler struct {
	svc UserService
}

func NewUserHandler(svc UserService) handler {
	return handler{
		svc: svc,
	}
}

func (h *handler) ListUsers(c *echo.Context) error {
	var params GetUserParams

	if err := c.Bind(&params); err != nil {
		return response.Error(c, apperror.BadRequest("invalid query params"))
	}

	users, total, err := h.svc.FetchAllUsers(c.Request().Context(), &params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.WithPagination(c, http.StatusOK, users, params.Page, params.Limit, int64(total))
}

func (h *handler) CreateUser(c *echo.Context) error {
	req, err := request.BindAndValidate[UserCreateRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.RegisterUser(c.Request().Context(), req); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusCreated, "user created successfully")
}

func (h *handler) GetUser(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, apperror.BadRequest("invalid user id"))
	}

	user, err := h.svc.FindUserByID(c.Request().Context(), &id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, user)
}

func (h *handler) UpdateUserAvatar(c *echo.Context) error {
	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		return response.Error(c, apperror.BadRequest("avatar is required"))
	}

	if err := validation.ValidateImage(fileHeader, maxAvatarSize, AllowedAvatarTypes); err != nil {
		return response.Error(c, err)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return response.Error(c, apperror.Internal("failed to open uploaded file"))
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("failed to close file", "error", err)
		}
	}()

	if err := h.svc.UpdateUserAvatar(c.Request().Context(), file, fileHeader); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusOK, "update avatar successfully")
}
