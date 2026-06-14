package seed

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func SeedUserAdmin(db *sqlx.DB) (bool, error) {
	const email = "admin@email.com"
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	const queryInsert = `
		INSERT INTO users (name, email, password, role_id, division_id, gender)
		VALUES ($1, $2, $3, (SELECT id FROM roles WHERE code = 'ADMIN'), $4, $5)
	`

	_, err = db.Exec(queryInsert, "Riyan Amanda", email, hashedPassword, 1, "MALE")
	if err != nil {
		return false, err
	}

	return true, nil
}
