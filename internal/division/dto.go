package division

import (
	"time"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"
)

type DivisionResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DivisionBrief struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GetDivisionParams struct {
	pagination.Params
	Search   string `query:"search"`
	SortBy   string `query:"sort_by"`
	SortType string `query:"sort_type"`
	IsActive *bool  `query:"is_active"`
}

type CreateDivisionRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type UpdateDivisionRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	IsActive *bool  `json:"is_active"`
}

func (p *GetDivisionParams) Normalize() {
	p.Params.Normalize()
}
