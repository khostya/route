package config

import "errors"

var (
	ErrTokensLengthIsNotValid    = errors.New("tokens length is not valid")
	ErrCapacityInGramInNotValid  = errors.New("capacity_in_gram is not valid")
	ErrPriceInRubIsNotValid      = errors.New("price_in_rub is not valid")
	ErrValueIsEmpty              = errors.New("value is empty")
	ErrWrappersConfigPathIsEmpty = errors.New("WRAPPERS_CONFIG_PATH is empty")
)
