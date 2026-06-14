package auth

import (
	"github.com/google/uuid"

	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type LoginResponse struct {
	User        CurrentUserResponse `json:"user"`
	AccessToken string              `json:"access_token"`
}

type DivisionBrief struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CurrentUserResponse struct {
	ID        uuid.UUID     `json:"id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	Role      user.UserRole `json:"role"`
	AvatarURL *string       `json:"avatar_url"`
	Division  DivisionBrief `json:"division"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type GoogleLoginRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}
