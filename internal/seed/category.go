package seed

import "github.com/jmoiron/sqlx"

func SeedCategory(db *sqlx.DB) error {
	const query = `
		INSERT INTO categories (name)
		VALUES ($1), ($2), ($3), ($4), ($5)
		ON CONFLICT DO NOTHING
	`

	_, err := db.Exec(query, "SIMRS", "Network", "Software", "Hardware", "Peripheral")
	if err != nil {
		return err
	}

	return nil
}
