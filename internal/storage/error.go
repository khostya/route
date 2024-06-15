package storage

import (
	"errors"
	"github.com/jackc/pgx"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrDuplicateOrderID = errors.New("duplicate order id")
)

func isDuplicateKeyError(err error) bool {
	var pgErr pgx.PgError
	ok := errors.As(err, &pgErr)
	if ok {
		// unique_violation = 23505
		return pgErr.Code == "23505"
	}
	return false
}
