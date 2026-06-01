package database

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

const uniqueViolationCode = "23505"
const foreignKeyViolationCode = "23503"

func IsUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && string(pqErr.Code) == uniqueViolationCode
}

func IsForeignKeyViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && string(pqErr.Code) == foreignKeyViolationCode
}

func CheckRowsAffected(result sql.Result, notFoundErr error) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return notFoundErr
	}

	return nil
}
