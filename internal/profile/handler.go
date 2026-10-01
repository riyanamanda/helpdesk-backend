package profile

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/validation"
)

type service interface {
	GetProfile(ctx context.Context) (*ProfileResponse, error)
	UpdateProfile(ctx context.Context, req *UpdateProfileRequest) error
	UpdateAvatar(ctx context.Context, file *storage.File) error
	SyncGoogle(ctx context.Context, req *SyncGoogleRequest) error
	RevokeGoogle(ctx context.Context) error
	UpdatePassword(ctx context.Context, req UpdatePasswordRequest) error
}

type Handler struct {
	svc service
}

func NewHandler(svc service) *Handler {
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
	req, err := httputil.BindAndValidate[UpdateProfileRequest](c)
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

	file := &storage.File{
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
	req, err := httputil.BindAndValidate[SyncGoogleRequest](c)
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
	req, err := httputil.BindAndValidate[UpdatePasswordRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.UpdatePassword(c.Request().Context(), *req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}
