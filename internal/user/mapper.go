package user

import (
	"github.com/riyanamanda/helpdesk-backend/internal/division"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/collection"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
)

func toUserResponse(u UserProjection, storageConfig config.Storage) UserResponse {
	var avatarURL *string
	var createdBy *UserBrief

	if u.AvatarKey != nil {
		url := utils.BuildPublicURL(storageConfig.PublicURL, storageConfig.Bucket, *u.AvatarKey)
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
		Role:      UserRole(u.Role),
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

func toUserResponses(users []UserProjection, storageConfig config.Storage) []UserResponse {
	return collection.MapSlice(users, func(u UserProjection) UserResponse {
		return toUserResponse(u, storageConfig)
	})
}
