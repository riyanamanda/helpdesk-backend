package feedback

func toFeedbackResponse(f FeedbackProjection) FeedbackResponse {
	return FeedbackResponse{
		ID:          f.ID,
		Title:       f.Title,
		Description: f.Description,
		Type:        f.Type,
		Status:      f.Status,
		CreatedBy: FeedbackUser{
			ID:   f.CreatedByID,
			Name: f.CreatedByName,
		},
		ReviewedBy: func() *FeedbackUser {
			if f.ReviewedByID == nil || f.ReviewedByName == nil {
				return nil
			}
			return &FeedbackUser{ID: *f.ReviewedByID, Name: *f.ReviewedByName}
		}(),
		ReviewedAt: f.ReviewedAt,
		CreatedAt:  f.CreatedAt,
		UpdatedAt:  f.UpdatedAt,
	}
}

func toFeedbackResponses(feedbacks []FeedbackProjection) []FeedbackResponse {
	result := make([]FeedbackResponse, len(feedbacks))
	for i, f := range feedbacks {
		result[i] = toFeedbackResponse(f)
	}
	return result
}
