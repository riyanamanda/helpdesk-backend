package category

import (
	"net/http"

	"github.com/labstack/echo/v5"
	apperrors "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
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

	category, err := h.svc.Create(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusCreated, category)
}
