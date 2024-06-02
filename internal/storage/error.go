package storage

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrDuplicateOrderID = fmt.Errorf("duplicate order id")
)
