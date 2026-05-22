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

	Search   string `query:"search"`
	SortBy   string `query:"sort_by"`
	SortType string `query:"sort_type"`
	IsActive *bool  `query:"is_active"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
}

type UpdateCategoryRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	IsActive *bool  `json:"is_active"`
}

func (p *GetCategoryParams) Normalize() {
	page, limit, _ := p.Params.Normalize()
	p.Page = page
	p.Limit = limit

	allowedSortBy := map[string]bool{
		"name":       true,
		"is_active":  true,
		"created_at": true,
	}

	if !allowedSortBy[p.SortBy] {
		p.SortBy = "created_at"
	}

	if p.SortType != "ASC" && p.SortType != "DESC" {
		p.SortType = "DESC"
	}
}
