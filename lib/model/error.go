package model

import (
	"errors"
	"regexp"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrorPersistence = errors.New("persistence error")
)

// postgres error code for unique constraint violations.
// see https://www.postgresql.org/docs/current/errcodes-appendix.html
const pgUniqueViolationCode = "23505"

// matches the "Key (col)=(value) already exists." detail postgres reports
// for a unique constraint violation, capturing the offending value.
var pgDuplicateKeyValueRegexp = regexp.MustCompile(`^Key \([^)]+\)=\(([^)]+)\) already exists\.$`)

// AsDuplicateKeyError reports whether err is a postgres unique constraint
// violation. When it is, it also returns the offending value, extracted
// from the driver error detail when possible.
func AsDuplicateKeyError(err error) (value string, ok bool) {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != pgUniqueViolationCode {
		return "", false
	}

	matches := pgDuplicateKeyValueRegexp.FindStringSubmatch(pgErr.Detail)
	if matches == nil {
		return "", true
	}
	return matches[1], true
}
