package service

import (
	"errors"
)

var (
	ErrOrderInPVZ                            = newError(errors.New("заказ находится в пвз"))
	ErrRefundPeriodHasExpired                = newError(errors.New("заказ не может быть возвращен более чем через два дня"))
	ErrOrderHasNotExpired                    = newError(errors.New("у заказа ещё не вышел срок хранения"))
	ErrOrderHasExpired                       = newError(errors.New("у заказа вышел срок хранения"))
	ErrOrderHasAlreadyBeenIssued             = newError(errors.New("заказ уже выдан"))
	ErrExtraIDsInTheRequest                  = newError(errors.New("в запросе присутствуют лишние id"))
	ErrExpIsNotValid                         = newError(errors.New("expiration date is not valid"))
	ErrOrdersBelongToDifferentUsers          = newError(errors.New("orders belong to different users"))
	ErrMustBeAtLeastOneOrder                 = newError(errors.New("must be at least one order"))
	ErrOrderWeightGreaterThanWrapperCapacity = newError(errors.New("order weight is greater than the wrapper capacity"))
)

type OrderServiceError struct {
	err error
}

func newError(err error) OrderServiceError {
	return OrderServiceError{err: err}
}

func (o OrderServiceError) Error() string {
	return o.err.Error()
}
