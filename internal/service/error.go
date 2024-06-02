package service

import (
	"errors"
)

var (
	ErrOrderInPVZ                = errors.New("заказ находится в пвз")
	ErrRefundPeriodHasExpired    = errors.New("заказ не может быть возвращен более чем через два дня")
	ErrOrderHasNotExpired        = errors.New("у заказа ещё не вышел срок хранения")
	ErrOrderHasExpired           = errors.New("у заказа вышел срок хранения")
	ErrOrderHasAlreadyBeenIssued = errors.New("заказ уже выдан")
	ErrExtraIDsInTheRequest      = errors.New("в запросе присутствуют лишние id")
)
