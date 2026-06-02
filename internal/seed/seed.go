package seed

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

func Run(db *sqlx.DB) error {
	categoryInserted, err := SeedCategory(db)
	if err != nil {
		return err
	}
	slog.Info("seed categories finished", "inserted", categoryInserted)

	divisionInserted, err := SeedDivision(db)
	if err != nil {
		return err
	}
	slog.Info("seed divisions finished", "inserted", divisionInserted)

	adminInserted, err := SeedUserAdmin(db)
	if err != nil {
		return err
	}
	if adminInserted {
		slog.Info("seed admin user finished", "inserted", true)
	} else {
		slog.Info("seed admin user finished", "inserted", false)
	}

	return nil
}
