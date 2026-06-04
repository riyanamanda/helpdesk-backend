package feedback

import "github.com/riyanamanda/helpdesk-backend/internal/shared/sliceutil"

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
	return sliceutil.Map(feedbacks, func(f FeedbackProjection) FeedbackResponse {
		return toFeedbackResponse(f)
	})
}
