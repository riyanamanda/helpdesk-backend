package user

import (
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
)

func toUserResponse(u UserProjection, storageConfig config.Storage) UserResponse {
	var avatarURL *string
	var createdBy *UserBrief

	if u.AvatarKey != nil {
		url := httputil.BuildPublicURL(storageConfig.Bucket, *u.AvatarKey)
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
		Gender:    u.Gender,
		Division: DivisionBrief{
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
	result := make([]UserResponse, len(users))
	for i, u := range users {
		result[i] = toUserResponse(u, storageConfig)
	}
	return result
}

func toUserBriefs(users []AssignableUserProjection) []UserBrief {
	result := make([]UserBrief, len(users))
	for i, u := range users {
		result[i] = UserBrief(u)
	}
	return result
}
