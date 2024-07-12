package config

import "errors"

var (
	ErrCapacityInGramIsNotValid  = errors.New("capacity_in_gram is not valid")
	ErrPriceInRubIsNotValid      = errors.New("price_in_rub is not valid")
	ErrWrappersConfigPathIsEmpty = errors.New("WRAPPERS_CONFIG_PATH is empty")
	ErrOutputConfigPathIsEmpty   = errors.New("OUTPUT_CONFIG_PATH is empty")
	ErrOutputFilterDoesNotExits  = errors.New("output filter does not exits")
	ErrKafkaConfigPathIsEmpty    = errors.New("KAFKA_CONFIG_PATH is empty")
)
