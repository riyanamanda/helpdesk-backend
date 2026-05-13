package seed

import "github.com/jmoiron/sqlx"

func SeedDivision(db *sqlx.DB) error {
	const query = `
		INSERT INTO divisions (name)
		VALUES ($1), ($2), ($3)
		ON CONFLICT DO NOTHING
	`

	_, err := db.Exec(query, "IT", "Rekam Medis", "Farmasi")
	if err != nil {
		return err
	}

	return nil
}
