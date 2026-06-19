package seed

import "github.com/jmoiron/sqlx"

func SeedPermission(db *sqlx.DB) (int64, error) {
	permissions := []struct {
		Code string
	}{
		// user module
		{"user:view"},
		{"user:create"},
		{"user:update"},
		{"user:delete"},

		// rbac module
		{"rbac:view"},
		{"rbac:create"},
		{"rbac:update"},
		{"rbac:delete"},

		// caregory module
		{"category:view"},
		{"category:create"},
		{"category:update"},
		{"category:delete"},

		// division module
		{"division:view"},
		{"division:create"},
		{"division:update"},
		{"division:delete"},

		// feedback module
		{"feedback:view"},
		{"feedback:create"},
		{"feedback:update"},
		{"feedback:delete"},

		// ticket module
		{"ticket:view"},
		{"ticket:create"},
		{"ticket:update"},
		{"ticket:delete"},
		{"ticket:assign"},
		{"ticket:priority"},
		{"ticket:resolution"},
		{"ticket:close"},

		// ihs module
		{"ihs:view"},
		{"ihs:update"},

		// antrian module
		{"antrian:view"},
		{"antrian:checkin"},
	}

	const query = `
		INSERT INTO permissions (code)
		VALUES ($1)
		ON CONFLICT DO NOTHING
	`

	var affected int64
	for _, p := range permissions {
		result, err := db.Exec(query, p.Code)
		if err != nil {
			return 0, err
		}

		n, err := result.RowsAffected()
		if err != nil {
			return 0, nil
		}

		affected += n
	}

	return affected, nil
}
