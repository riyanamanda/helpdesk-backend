package auth

import (
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func toCurrentUserResponse(u user.UserProjection, storageConfig config.Storage) CurrentUserResponse {
	var avatarURL *string

	if u.AvatarKey != nil {
		url := httputil.BuildPublicURL(storageConfig.Bucket, *u.AvatarKey)
		avatarURL = &url
	}

	return CurrentUserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Role:  u.Role,
		Division: DivisionBrief{
			ID:   u.DivisionID,
			Name: u.DivisionName,
		},
		AvatarURL: avatarURL,
	}
}
