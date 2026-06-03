package profile

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/request"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/upload"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/validation"
)

type Handler struct {
	svc ProfileService
}

func NewProfileHandler(svc ProfileService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) GetProfile(c *echo.Context) error {
	profile, err := h.svc.GetProfile(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, profile)
}

func (h *Handler) UpdateProfile(c *echo.Context) error {
	req, err := request.BindAndValidate[UpdateProfileRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.UpdateProfile(c.Request().Context(), req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) UpdateAvatar(c *echo.Context) error {
	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		return response.Error(c, apperr.BadRequest("avatar is required"))
	}

	if err := validation.ValidateImage(fileHeader, maxAvatarSize, allowedAvatarTypes); err != nil {
		return response.Error(c, err)
	}

	f, err := fileHeader.Open()
	if err != nil {
		return response.Error(c, apperr.Internal())
	}
	defer f.Close()

	file := &upload.File{
		Content:     f,
		Filename:    fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
		Size:        fileHeader.Size,
	}

	if err := h.svc.UpdateAvatar(c.Request().Context(), file); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) SyncGoogle(c *echo.Context) error {
	req, err := request.BindAndValidate[SyncGoogleRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.SyncGoogle(c.Request().Context(), req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) RevokeGoogle(c *echo.Context) error {
	if err := h.svc.RevokeGoogle(c.Request().Context()); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) UpdatePassword(c *echo.Context) error {
	req, err := request.BindAndValidate[UpdatePasswordRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.UpdatePassword(c.Request().Context(), *req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}
