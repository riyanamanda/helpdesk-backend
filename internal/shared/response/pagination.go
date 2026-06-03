package response

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type PaginationResponse struct {
	Data       any            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

type PaginationMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"total_page"`
}

func Paginated(
	c *echo.Context,
	data any,
	page int,
	limit int,
	total int64,
) error {
	totalPage := (total + int64(limit) - 1) / int64(limit)

	return c.JSON(http.StatusOK, PaginationResponse{
		Data: data,
		Pagination: PaginationMeta{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: totalPage,
		},
	})
}
