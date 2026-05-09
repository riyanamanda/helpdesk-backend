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

type GetDivisionParams struct {
	pagination.Params
}

type CreateDivisionRequest struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
}

type UpdateDivisionRequest struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
}
