package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"
)

type UserResponse struct {
	ID        uuid.UUID     `json:"id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	GoogleID  *string       `json:"google_id"`
	AvatarURL *string       `json:"avatar_url"`
	Phone     *string       `json:"phone"`
	Role      UserRole      `json:"role"`
	Gender    string        `json:"gender"`
	Division  DivisionBrief `json:"division"`
	IsActive  bool          `json:"is_active"`
	CreatedBy *UserBrief    `json:"created_by"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type DivisionBrief struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserBrief struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type GetUserParams struct {
	pagination.Params
	Search   string   `query:"search"`
	SortBy   string   `query:"sort_by"`
	SortType string   `query:"sort_type"`
	IsActive *bool    `query:"is_active"`
	Role     UserRole `query:"role"`
	Division *int64   `query:"division"`
}

type UserCreateRequest struct {
	Name     string   `json:"name" validate:"required,min=3,max=20"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	Role     UserRole `json:"role" validate:"required,oneof=ADMIN EMPLOYEE"`
	Division int64    `json:"division" validate:"required,gt=0"`
	Gender   string   `json:"gender" validate:"required"`
}

type UserUpdateRequest struct {
	Name     string   `json:"name" validate:"required,min=3,max=20"`
	Email    string   `json:"email" validate:"required,email"`
	Role     UserRole `json:"role" validate:"required,oneof=ADMIN EMPLOYEE"`
	Division int64    `json:"division" validate:"required,gt=0"`
	Gender   string   `json:"gender" validate:"required"`
	IsActive *bool    `json:"is_active" validate:"required"`
}

type UserUpdatePassword struct {
	Password string `json:"password" validate:"required,min=8"`
}

func (p *GetUserParams) Normalize() {
	p.Params.Normalize()
}
