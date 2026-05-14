package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"
)

type UserResponse struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	AvatarURL  *string    `json:"avatar_url"`
	Phone      *string    `json:"phone"`
	Role       UserRole   `json:"role"`
	DivisionID int64      `json:"division_id"`
	IsActive   bool       `json:"is_active"`
	CreatedBy  *uuid.UUID `json:"created_by"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type GetUserParams struct {
	pagination.Params
}

type UserCreateRequest struct {
	Name       string   `json:"name" validate:"required,min=3,max=20"`
	Email      string   `json:"email" validate:"required,email"`
	Password   string   `json:"password" validate:"required,min=8"`
	Role       UserRole `json:"role" validate:"required,oneof=ADMIN EMPLOYEE"`
	DivisionID int64    `json:"division_id" validate:"gt=0"`
}
