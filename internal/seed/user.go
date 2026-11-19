package seed

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/passwordtoken"
)

func SeedUserAdmin(db *sqlx.DB) (bool, error) {
	const email = "ryant.n92@gmail.com"
	var exists bool

	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE LOWER(email) = LOWER($1)
		)
	`

	err := db.Get(&exists, query, email)
	if err != nil {
		return false, err
	}

	if exists {
		slog.Info("user admin already exists")
		return false, nil
	}

	random, err := passwordtoken.New()
	if err != nil {
		return false, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(random), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	const queryInsert = `
		INSERT INTO users (name, email, password, role_id, division_id, gender)
		VALUES ($1, $2, $3, (SELECT id FROM roles WHERE code = 'SUPERADMIN'), $4, $5)
	`

	_, err = db.Exec(queryInsert, "Riyan Amanda", email, hashedPassword, 1, "MALE")
	if err != nil {
		return false, err
	}

	slog.Info("user admin created with a random password, use forgot password to set it")

	return true, nil
}
