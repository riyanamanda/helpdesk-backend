package profile

import (
	"time"

	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type DivisionBrief struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ProfileResponse struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	GoogleID  *string         `json:"google_id"`
	AvatarURL *string         `json:"avatar_url"`
	Phone     *string         `json:"phone"`
	Role      user.UserRole   `json:"role"`
	Gender    string          `json:"gender"`
	Division  DivisionBrief   `json:"division"`
	IsActive  bool            `json:"is_active"`
	CreatedBy *user.UserBrief `json:"created_by"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type UpdateProfileRequest struct {
	Name   string  `json:"name" validate:"required,min=3,max=50"`
	Phone  *string `json:"phone" validate:"omitempty,min=10,max=15"`
	Gender string  `json:"gender" validate:"required"`
}

type SyncGoogleRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}
