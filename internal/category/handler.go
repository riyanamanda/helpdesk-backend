package category

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	apperrors "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
	sharedutils "github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
)

type handler struct {
	svc CategoryService
}

func NewCategoryHandler(svc CategoryService) *handler {
	return &handler{
		svc: svc,
	}
}

func (h *handler) ListCategories(c *echo.Context) error {
	var params GetCategoryParams

	if err := c.Bind(&params); err != nil {
		return response.Error(c, apperrors.BadRequest("invalid query params"))
	}

	page, limit, _ := params.Normalize()

	categories, total, err := h.svc.GetCategories(c.Request().Context(), &params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.WithPagination(c, http.StatusOK, categories, page, limit, total)
}

func (h *handler) Create(c *echo.Context) error {
	var req CreateCategoryRequest

	if err := c.Bind(&req); err != nil {
		return response.Error(c, apperrors.BadRequest("invalid request format"))
	}
	req.Name = strings.TrimSpace(req.Name)

	category, err := h.svc.Create(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusCreated, category)
}

func (h *handler) GetByID(c *echo.Context) error {
	id, err := sharedutils.ParsePositiveInt64PathParam(c, "id", "category")
	if err != nil {
		return response.Error(c, err)
	}

	category, err := h.svc.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, category)
}

func (h *handler) Update(c *echo.Context) error {
	id, err := sharedutils.ParsePositiveInt64PathParam(c, "id", "category")
	if err != nil {
		return response.Error(c, err)
	}

	var req UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, apperrors.BadRequest("invalid request format"))
	}
	req.Name = strings.TrimSpace(req.Name)

	category, err := h.svc.Update(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, category)
}

func (h *handler) Delete(c *echo.Context) error {
	id, err := sharedutils.ParsePositiveInt64PathParam(c, "id", "category")
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return response.Error(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
