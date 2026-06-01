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

type DivisionOptionResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GetDivisionParams struct {
	pagination.Params

	IsActive *bool `query:"is_active"`
}

type DivisionCreateRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type DivisionUpdateRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	IsActive *bool  `json:"is_active"`
}

func (p *GetDivisionParams) Normalize() {
	p.Params.Normalize()
}
