package psql

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	UNIQUE_VIOLATION_CODE = "23505"
)

func checkUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == UNIQUE_VIOLATION_CODE {
			return true
		}
	}

	return false
}
