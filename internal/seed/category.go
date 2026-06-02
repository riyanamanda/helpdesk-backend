package seed

import "github.com/jmoiron/sqlx"

func SeedCategory(db *sqlx.DB) (int64, error) {
	const query = `
		INSERT INTO categories (name)
		VALUES ($1), ($2), ($3), ($4), ($5)
		ON CONFLICT DO NOTHING
	`

	result, err := db.Exec(query, "SIMRS", "Network", "Software", "Hardware", "Peripheral")
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
