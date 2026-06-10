package user_device

import (
	"time"

	"github.com/google/uuid"
)

type UserDevice struct {
	ID          int64     `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	FcmToken    string    `db:"fcm_token"`
	LastSeenAt  time.Time `db:"last_seen_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
