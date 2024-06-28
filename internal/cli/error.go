package cli

import (
	"errors"
)

var (
	ErrIdIsEmpty            = errors.New("id is empty")
	ErrUserIsEmpty          = errors.New("user is empty")
	ErrExpIsEmpty           = errors.New("exp is empty")
	ErrPageIsNotValid       = errors.New("page is not valid")
	ErrSizeIsNotValid       = errors.New("size is not valid")
	ErrWrapperIsNotValid    = errors.New("wrapper is not valid")
	ErrWeightInKgInNotValid = errors.New("weight_in_kg is not valid")
	ErrPriceInRubIsNotValid = errors.New("price_in_rub is not valid")
	ErrExit                 = errors.New(exit)
)
