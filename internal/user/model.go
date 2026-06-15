package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID  `db:"id"`
	Name       string     `db:"name"`
	Email      string     `db:"email"`
	Password   string     `db:"password" json:"-"`
	GoogleID   *string    `db:"google_id"`
	AvatarKey  *string    `db:"avatar_key"`
	Phone      *string    `db:"phone"`
	RoleID     int64      `db:"role_id"`
	Gender     string     `db:"gender"`
	DivisionID int64      `db:"division_id"`
	IsActive   bool       `db:"is_active"`
	CreatedBy  *uuid.UUID `db:"created_by"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}
