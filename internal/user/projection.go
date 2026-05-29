package user

import (
	"time"

	"github.com/google/uuid"
)

type UserProjection struct {
	ID            uuid.UUID  `db:"id"`
	Name          string     `db:"name"`
	Email         string     `db:"email"`
	Password      string     `db:"password"`
	GoogleID      *string    `db:"google_id"`
	AvatarKey     *string    `db:"avatar_key"`
	Phone         *string    `db:"phone"`
	Role          UserRole   `db:"role"`
	Gender        string     `db:"gender"`
	DivisionID    int64      `db:"division_id"`
	DivisionName  string     `db:"division_name"`
	IsActive      bool       `db:"is_active"`
	CreatedByID   *uuid.UUID `db:"created_by_id"`
	CreatedByName *string    `db:"created_by_name"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
}
