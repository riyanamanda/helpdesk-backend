package ihs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PatientRepository interface {
	GetPatients(ctx context.Context, params GetPatientParams) ([]Patient, int64, error)
	GetPatientDetail(ctx context.Context, NORM string) (*PatientDetailProjection, error)
	UpdatePatientMethod(ctx context.Context, NORM string) error
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

func (r *repository) GetPatientDetail(ctx context.Context, NORM string) (*PatientDetailProjection, error) {
	var patient PatientDetailProjection

	const query = `
		SELECT
			p.NORM as norm,
			p.NAMA as name,
			wil.DESKRIPSI as birth_place,
			p.TANGGAL_LAHIR as birth_date,
			p.JENIS_KELAMIN as gender,
			k_ref.DESKRIPSI as marital_status,
			w_ref.DESKRIPSI as citizenship,
			p.STATUS as status,
			ktp.NOMOR as identity_number,
			ktp.ALAMAT as address,
			ktp.RT as rt,
			ktp.RW as rw,
			prov.DESKRIPSI as province,
			ct.DESKRIPSI as city,
			ds.DESKRIPSI as district,
			sds.DESKRIPSI as sub_district
		FROM master.pasien p
		LEFT JOIN master.kartu_identitas_pasien ktp
			ON ktp.NORM = p.NORM and ktp.JENIS = 1
		LEFT JOIN master.wilayah wil
			ON wil.id = p.TEMPAT_LAHIR
		JOIN master.referensi w_ref
			ON w_ref.id = p.KEWARGANEGARAAN and w_ref.JENIS = 177
		JOIN master.referensi k_ref
			ON k_ref.id = p.STATUS_PERKAWINAN and k_ref.JENIS = 5
		LEFT JOIN master.wilayah prov
			ON prov.ID = LEFT(ktp.WILAYAH, 2)
		LEFT JOIN master.wilayah ct
			ON ct.ID = LEFT(ktp.WILAYAH, 4)
		LEFT JOIN master.wilayah ds
			ON ds.ID = LEFT(ktp.WILAYAH, 6)
		LEFT JOIN master.wilayah sds
			ON sds.ID = ktp.WILAYAH
		WHERE p.NORM = ?
	`

	if err := r.db.GetContext(ctx, &patient, query, NORM); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPatientNotFound
		}
		return nil, err
	}

	return &patient, nil
}

func (r *repository) UpdatePatientMethod(ctx context.Context, NORM string) error {
	const query = `
		UPDATE ` + "`kemkes-ihs`" + `.patient
		SET httpRequest = 'POST'
		WHERE refId = ?
		AND id IS NULL
		AND statusRequest = 0
	`

	result, err := r.db.ExecContext(ctx, query, NORM)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrPatientNotFound
	}

	return nil
}
