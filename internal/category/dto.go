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

type CategoryOptionResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GetCategoryParams struct {
	pagination.Params

	IsActive *bool `query:"is_active"`
}

type CategoryCreateRequest struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
}

type CategoryUpdateRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	IsActive *bool  `json:"is_active"`
}

func (p *GetCategoryParams) Normalize() {
	p.Params.Normalize()
}
