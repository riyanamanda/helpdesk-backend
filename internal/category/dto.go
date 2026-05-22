package category

import (
	"time"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"
)

type CategoryResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryBrief struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GetCategoryParams struct {
	pagination.Params
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
}

type UpdateCategoryRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	IsActive *bool  `json:"is_active"`
}
