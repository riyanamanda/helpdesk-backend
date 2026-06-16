package seed

import "github.com/jmoiron/sqlx"

func SeedDivision(db *sqlx.DB) (int64, error) {
	const query = `
		INSERT INTO divisions (name)
		VALUES ($1), ($2), ($3), ($4), ($5), ($6), ($7)
		ON CONFLICT DO NOTHING
	`

	result, err := db.Exec(query, "IT")
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
