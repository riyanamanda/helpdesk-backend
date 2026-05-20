package auth

import (
	"github.com/google/uuid"
	"github.com/riyanamanda/helpdesk-backend/internal/division"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User        CurrentUserResponse `json:"user"`
	AccessToken string              `json:"access_token"`
}

type CurrentUserResponse struct {
	ID        uuid.UUID              `json:"id"`
	Name      string                 `json:"name"`
	Email     string                 `json:"email"`
	Role      user.UserRole          `json:"role"`
	AvatarURL *string                `json:"avatar_url"`
	Division  division.DivisionBrief `json:"division"`
}
