package ihs

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PatientRepository interface {
	GetPatients(ctx context.Context, params GetPatientParams) ([]Patient, int64, error)
}

type repository struct {
	db *sqlx.DB
}

func NewPatientRepository(db *sqlx.DB) PatientRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetPatients(ctx context.Context, params GetPatientParams) ([]Patient, int64, error) {
	var (
		patients []Patient
		total    int64
	)

	where, args := buildPatientWhere(params)

	queryTotal := fmt.Sprintf("SELECT COUNT(*) FROM `kemkes-ihs`.patient ip JOIN master.pasien p ON ip.refId = p.NORM %s", where)
	if err := r.db.GetContext(ctx, &total, queryTotal, args...); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)

	col, dir := buildPatientSort(params)
	query := fmt.Sprintf(patientSelectBase+`
	%s
	ORDER BY %s %s
	LIMIT ? OFFSET ?
	`, where, col, dir)

	if err := r.db.SelectContext(ctx, &patients, query, args...); err != nil {
		return nil, 0, err
	}

	return patients, total, nil
}
