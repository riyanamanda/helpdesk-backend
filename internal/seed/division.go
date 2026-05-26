package seed

import "github.com/jmoiron/sqlx"

func SeedDivision(db *sqlx.DB) error {
	const query = `
		INSERT INTO divisions (name)
		VALUES ($1), ($2), ($3), ($4), ($5), ($6), ($7)
		ON CONFLICT DO NOTHING
	`

	_, err := db.Exec(query, "IT", "Rekam Medis", "Pharmachy", "IGD", "Physiology", "Laboratory", "Radiology")
	if err != nil {
		return err
	}

	return nil
}
