package category

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/response"
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
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	params := ListCategoriesParams{
		Limit:  limit,
		Offset: offset,
	}

	result, err := h.svc.GetCategories(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("internal server error"))
	}

	meta := &response.Meta{
		Limit:  params.Limit,
		Offset: params.Offset,
		Total:  result.Total,
	}

	return c.JSON(http.StatusOK,
		response.Success(toCategoryResponses(result.Data), meta),
	)
}

func (h *handler) Create(c *echo.Context) error {
	var req CreateCategoryRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request format"))
	}

	category, err := h.svc.Create(c.Request().Context(), &req)
	if err != nil {
		if errors.Is(err, ErrInvalidCategory) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		if errors.Is(err, ErrCategoryAlreadyExists) {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": category,
	})
}
