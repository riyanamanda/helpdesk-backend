package notification

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID            int64                     `db:"id"`
	UserID        uuid.UUID                 `db:"user_id"`
	Type          NotificationType          `db:"type"`
	ReferenceType NotificationReferenceType `db:"reference_type"`
	ReferenceID   int64                     `db:"reference_id"`
	Metadata      string                    `db:"metadata"`
	IsRead        bool                      `db:"is_read"`
	ReadAt        *time.Time                `db:"read_at"`
	CreatedAt     time.Time                 `db:"created_at"`
}
