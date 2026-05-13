package database

import (
	"errors"

	"github.com/lib/pq"
)

const uniqueViolationCode = "23505"

func IsUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && string(pqErr.Code) == uniqueViolationCode
}