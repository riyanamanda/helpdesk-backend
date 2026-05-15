package user

import (
	"github.com/riyanamanda/helpdesk-backend/internal/division"
	"github.com/riyanamanda/helpdesk-backend/internal/storage"
)

func toUserResponse(u UserProjection, storage storage.Storage) UserResponse {
	var avatarURL *string
	var createdBy *UserBrief

	if u.AvatarKey != nil {
		url := storage.GetURL(*u.AvatarKey)
		avatarURL = &url
	}

	if u.CreatedByID != nil && u.CreatedByName != nil {
		createdBy = &UserBrief{
			ID:   *u.CreatedByID,
			Name: *u.CreatedByName,
		}
	}

	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		GoogleID:  u.GoogleID,
		AvatarURL: avatarURL,
		Phone:     u.Phone,
		Role:      u.Role,
		Division: division.DivisionBrief{
			ID:   u.DivisionID,
			Name: u.DivisionName,
		},
		IsActive:  u.IsActive,
		CreatedBy: createdBy,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func toUserResponses(users []UserProjection, storage storage.Storage) []UserResponse {
	result := make([]UserResponse, len(users))
	for i, u := range users {
		result[i] = toUserResponse(u, storage)
	}
	return result
}
