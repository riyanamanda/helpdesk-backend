package seed

import "github.com/jmoiron/sqlx"

func Run(db *sqlx.DB) error {
	if err := SeedCategory(db); err != nil {
		return err
	}

	if err := SeedDivision(db); err != nil {
		return err
	}

	if err := SeedUserAdmin(db); err != nil {
		return err
	}

	return nil
}
