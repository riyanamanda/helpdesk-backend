package feedback

type (
	FeedbackType   string
	FeedbackStatus string
)

const (
	FeatureRequest FeedbackType = "FEATURE_REQUEST"
	Improvement    FeedbackType = "IMPROVEMENT"
	BugReport      FeedbackType = "BUG_REPORT"

	FeedbackStatusOpen      FeedbackStatus = "OPEN"
	FeedbackStatusInReview  FeedbackStatus = "IN_REVIEW"
	FeedbackStatusAccepted  FeedbackStatus = "ACCEPTED"
	FeedbackStatusRejected  FeedbackStatus = "REJECTED"
	FeedbackStatusDelivered FeedbackStatus = "DELIVERED"
)
