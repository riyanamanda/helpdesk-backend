package notification

import (
	"encoding/json"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/sliceutil"
)

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
	return sliceutil.Map(notifications, func(n Notification) NotificationResponse {
		return toNotificationResponse(n)
	})
}
