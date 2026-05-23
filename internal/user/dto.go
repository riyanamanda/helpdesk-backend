package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/division"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"
)

type UserResponse struct {
	ID        uuid.UUID              `json:"id"`
	Name      string                 `json:"name"`
	Email     string                 `json:"email"`
	GoogleID  *string                `json:"google_id"`
	AvatarURL *string                `json:"avatar_url"`
	Phone     *string                `json:"phone"`
	Role      UserRole               `json:"role"`
	Division  division.DivisionBrief `json:"division"`
	IsActive  bool                   `json:"is_active"`
	CreatedBy *UserBrief             `json:"created_by"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type UserBrief struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type GetUserParams struct {
	pagination.Params
}

type UserCreateRequest struct {
	Name       string   `json:"name" validate:"required,min=3,max=20"`
	Email      string   `json:"email" validate:"required,email"`
	Password   string   `json:"password" validate:"required,min=8"`
	Role       UserRole `json:"role" validate:"required,oneof=ADMIN EMPLOYEE"`
	DivisionID int64    `json:"division_id" validate:"required,gt=0"`
}
