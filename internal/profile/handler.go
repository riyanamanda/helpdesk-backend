package profile

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperror"
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

	return response.Success(c, http.StatusOK, profile)
}

func (h *Handler) UpdateProfile(c *echo.Context) error {
	req, err := request.BindAndValidate[UpdateProfileRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.UpdateProfile(c.Request().Context(), req); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusOK, "profile updated successfully")
}

func (h *Handler) UpdateAvatar(c *echo.Context) error {
	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		return response.Error(c, apperror.BadRequest("avatar is required"))
	}

	if err := validation.ValidateImage(fileHeader, maxAvatarSize, allowedAvatarTypes); err != nil {
		return response.Error(c, err)
	}

	f, err := fileHeader.Open()
	if err != nil {
		return response.Error(c, apperror.Internal("failed to open uploaded file"))
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

	return response.Message(c, http.StatusOK, "avatar updated successfully")
}

func (h *Handler) SyncGoogle(c *echo.Context) error {
	req, err := request.BindAndValidate[SyncGoogleRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.SyncGoogle(c.Request().Context(), req); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusOK, "google account linked successfully")
}

func (h *Handler) RevokeGoogle(c *echo.Context) error {
	if err := h.svc.RevokeGoogle(c.Request().Context()); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusOK, "google account unlinked successfully")
}

func (h *Handler) UpdatePassword(c *echo.Context) error {
	req, err := request.BindAndValidate[UpdatePasswordRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.UpdatePassword(c.Request().Context(), *req); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusOK, "Password updated successfully")
}
