package feedback

type (
	FeedbackType   string
	FeedbackStatus string
)

const (
	FeatureRequest FeedbackType = "FEATURE_REQUEST"
	Improvement    FeedbackType = "IMPROVEMENT"
	BugReport      FeedbackType = "BUG_REPORT"

	Open      FeedbackStatus = "OPEN"
	InReview  FeedbackStatus = "IN_REVIEW"
	Accepted  FeedbackStatus = "ACCEPTED"
	Rejected  FeedbackStatus = "REJECTED"
	Delivered FeedbackStatus = "DELIVERED"
)
