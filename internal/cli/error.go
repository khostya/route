package cli

import (
	"errors"
	"fmt"
)

var (
	ErrIdIsEmpty      = errors.New("id is empty")
	ErrUserIsEmpty    = errors.New("user is empty")
	ErrExpIsEmpty     = errors.New("exp is empty")
	ErrPageIsNotValid = errors.New("page is not valid")
	ErrSizeIsNotValid = errors.New("size is not valid")
	ErrExit           = fmt.Errorf(exit)
)
