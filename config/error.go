package config

import "errors"

var (
	ErrCapacityInGramInNotValid  = errors.New("capacity_in_gram is not valid")
	ErrPriceInRubIsNotValid      = errors.New("price_in_rub is not valid")
	ErrWrappersConfigPathIsEmpty = errors.New("WRAPPERS_CONFIG_PATH is empty")
)
