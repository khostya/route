package storage

import (
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrDuplicateOrderID = errors.New("duplicate order id")
)

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	ok := errors.As(err, &pgErr)
	if ok {
		return pgErr.Code == pgerrcode.UniqueViolation
	}
	return false
}
