package notification

import "encoding/json"

func toNotificationResponse(n Notification) NotificationResponse {
	return NotificationResponse{
		ID:            n.ID,
		UserID:        n.UserID,
		Type:          string(n.Type),
		ReferenceType: n.ReferenceType,
		ReferenceID:   n.ReferenceID,
		Metadata:      json.RawMessage(n.Metadata),
		IsRead:        n.IsRead,
		ReadAt:        n.ReadAt,
		CreatedAt:     n.CreatedAt,
	}
}

func toNotificationResponses(notifications []Notification) []NotificationResponse {
	result := make([]NotificationResponse, len(notifications))
	for i, n := range notifications {
		result[i] = toNotificationResponse(n)
	}
	return result
}
