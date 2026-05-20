package auth

import (
	"github.com/riyanamanda/helpdesk-backend/internal/division"
	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func toCurrentUserResponse(u user.UserProjection, storageConfig config.Storage) CurrentUserResponse {
	var avatarURL *string
	if u.AvatarKey != nil {
		url := utils.BuildPublicURL(storageConfig.PublicURL, storageConfig.Bucket, *u.AvatarKey)
		avatarURL = &url
	}

	return CurrentUserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Role:  u.Role,
		Division: division.DivisionBrief{
			ID:   u.DivisionID,
			Name: u.DivisionName,
		},
		AvatarURL: avatarURL,
	}
}
