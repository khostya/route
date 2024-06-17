package storage

import (
	"errors"
	"github.com/jackc/pgconn"
)

var (
	ErrNotFound                      = errors.New("not found")
	ErrDuplicateOrderID              = errors.New("duplicate order id")
	ErrListWithHashesDifferentLength = errors.New("different length")
)

func isDuplicateKeyError(err error) bool {
	var pgErr = new(pgconn.PgError)

	ok := errors.As(err, &pgErr)
	if ok {
		// unique_violation = 23505
		return pgErr.Code == "23505"
	}
	return false
}
