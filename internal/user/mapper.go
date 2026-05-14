package user

import "github.com/riyanamanda/helpdesk-backend/internal/storage"

func toUserResponse(u User, storage storage.Storage) UserResponse {
	var avatarURL *string

	if u.AvatarKey != nil {
		url := storage.GetURL(*u.AvatarKey)
		avatarURL = &url
	}

	return UserResponse{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		AvatarURL:  avatarURL,
		Phone:      u.Phone,
		Role:       u.Role,
		DivisionID: u.DivisionID,
		IsActive:   u.IsActive,
		CreatedBy:  u.CreatedBy,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func toUserResponses(users []User, storage storage.Storage) []UserResponse {
	result := make([]UserResponse, 0, len(users))
	for _, u := range users {
		result = append(result, toUserResponse(u, storage))
	}
	return result
}
