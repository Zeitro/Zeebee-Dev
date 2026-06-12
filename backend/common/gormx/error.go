package gormx

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// pgUniqueViolationCode is the PostgreSQL SQLSTATE for unique_violation.
const pgUniqueViolationCode = "23505"

// IsUniqueViolation reports whether err is (or wraps) a PostgreSQL unique
// constraint violation (SQLSTATE 23505). It traverses wrapped errors via
// errors.As, so it works with errors returned through errorx.WithStack.
func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == pgUniqueViolationCode
	}
	return false
}
