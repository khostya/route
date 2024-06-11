package storage

import (
	"errors"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrDuplicateOrderID = errors.New("duplicate order id")
)
