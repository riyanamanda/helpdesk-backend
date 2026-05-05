package category

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
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

	categories, err := h.svc.GetCategories(c.Request().Context(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": toCategoryResponses(categories),
		"meta": map[string]int{
			"limit":  params.Limit,
			"offset": params.Offset,
		},
	})
}
