package antrian

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAntrian(ctx context.Context, params GetAntrianParams) ([]Antrian, int64, error) {
	var (
		antrian []Antrian
		total   int64
	)

	where, args := buildAntrianWhere(params)

	queryTotal := fmt.Sprintf("SELECT COUNT(*) FROM regonline.reservasi r LEFT JOIN regonline.dokter d ON d.KODE = r.DOKTER %s", where)
	if err := r.db.GetContext(ctx, &total, queryTotal, args...); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)

	query := fmt.Sprintf(antrianSelectBase+`
		%s
		ORDER BY r.STATUS DESC, r.NO ASC
		LIMIT ? OFFSET ?
		`, where)

	if err := r.db.SelectContext(ctx, &antrian, query, args...); err != nil {
		return nil, 0, err
	}

	return antrian, total, nil
}
