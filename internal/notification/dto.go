package notification

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID            int64                     `json:"id"`
	UserID        uuid.UUID                 `json:"user_id"`
	Type          string                    `json:"type"`
	ReferenceType NotificationReferenceType `json:"reference_type"`
	ReferenceID   int64                     `json:"reference_id"`
	Metadata      json.RawMessage           `json:"metadata"`
	IsRead        bool                      `json:"is_read"`
	ReadAt        *time.Time                `json:"read_at"`
	CreatedAt     time.Time                 `json:"created_at"`
}

type NotificationMetadata struct {
	ActorName string `json:"actor_name"`
	Status    string `json:"status,omitempty"`
}

type UnreadCountResponse struct {
	Count int64 `json:"count"`
}
