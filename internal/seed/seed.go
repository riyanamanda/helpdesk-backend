package seed

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

func Run(db *sqlx.DB) error {
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

	permissionsInserted, err := SeedPermission(db)
	if err != nil {
		return err
	}
	slog.Info("seed permissions finished", "inserted", permissionsInserted)

	rolePermissionInserted, err := SeedRolePermission(db)
	if err != nil {
		return err
	}
	slog.Info("seed role permissions finished", "inserted", rolePermissionInserted)

	return nil
}
