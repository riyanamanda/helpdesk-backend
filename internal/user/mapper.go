package user

func toUserResponse(u User) UserResponse {
	return UserResponse{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		AvatarKey:  u.AvatarKey,
		Phone:      u.Phone,
		DivisionID: u.DivisionID,
		IsActive:   u.IsActive,
		CreatedBy:  u.CreatedBy,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func toUserResponses(users []User) []UserResponse {
	result := make([]UserResponse, 0, len(users))
	for _, u := range users {
		result = append(result, toUserResponse(u))
	}
	return result
}
