package profile

import (
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func toProfileResponse(p user.UserProjection, storageConfig config.Storage) ProfileResponse {
	var avatarURL *string
	var createdBy *user.UserBrief

	if p.AvatarKey != nil {
		url := httputil.BuildPublicURL(storageConfig.PublicURL, storageConfig.Bucket, *p.AvatarKey)
		avatarURL = &url
	}

	if p.CreatedByID != nil && p.CreatedByName != nil {
		createdBy = &user.UserBrief{
			ID:   *p.CreatedByID,
			Name: *p.CreatedByName,
		}
	}

	return ProfileResponse{
		ID:        p.ID,
		Name:      p.Name,
		Email:     p.Email,
		GoogleID:  p.GoogleID,
		AvatarURL: avatarURL,
		Phone:     p.Phone,
		Role:      p.Role,
		Gender:    p.Gender,
		Division: DivisionBrief{
			ID:   p.DivisionID,
			Name: p.DivisionName,
		},
		IsActive:  p.IsActive,
		CreatedBy: createdBy,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
