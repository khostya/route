//go:generate mockgen -source ./mocks/order.go -destination=./mocks/mock_order.go -package=mock_service
package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"homework/internal/dto"
	"homework/internal/model"
	"homework/internal/model/wrapper"
	hash2 "homework/pkg/hash"
	"time"
)

const (
	refundPeriod = time.Hour * 24 * 2
)

type (
	orderStorage interface {
		ListUserOrders(ctx context.Context, id string, count uint, status model.Status) ([]model.Order, error)
		AddOrder(ctx context.Context, order model.Order, hash string) error
		ListOrdersByIds(ctx context.Context, ids []string, status model.Status) ([]model.Order, error)
		UpdateStatus(ctx context.Context, ids dto.IdsWithHashes, status model.Status) error
		GetOrderById(ctx context.Context, id string) (model.Order, error)
		DeleteOrder(ctx context.Context, id string) error
		RefundedOrders(ctx context.Context, get dto.PageParam) ([]model.Order, error)
		ListOrders(ctx context.Context, get dto.ListOrdersParam) ([]model.Order, error)
	}

	wrapperStorage interface {
		AddWrapper(ctx context.Context, order wrapper.Wrapper, orderID string) error
		Delete(ctx context.Context, orderID string) error
	}

	transactionManager interface {
		RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error
		Unwrap(err error) error
	}

	Deps struct {
		Storage            orderStorage
		TransactionManager transactionManager
		WrapperStorage     wrapperStorage
	}

	Order struct {
		orderStorage       orderStorage
		transactionManager transactionManager
		wrapperStorage     wrapperStorage
	}
)

func NewOrder(d Deps) Order {
	return Order{
		orderStorage:       d.Storage,
		transactionManager: d.TransactionManager,
		wrapperStorage:     d.WrapperStorage,
	}
}

func (o *Order) Deliver(ctx context.Context, order dto.DeliverOrderParam) error {
	if order.ExpirationDate.Before(time.Now()) {
		return ErrExpIsNotValid
	}
	if order.Wrapper != nil && !order.Wrapper.WillFitGram(order.WeightInGram) {
		message := fmt.Sprintf("capacity_in_gram = %v", order.Wrapper.GetCapacityInGram())
		return errors.Wrap(ErrOrderWeightGreaterThanWrapperCapacity, message)
	}

	wrapperPriceInRub := wrapper.PriceInRub(decimal.NewFromInt(0))
	if order.Wrapper != nil {
		wrapperPriceInRub = order.Wrapper.GetPriceInRub()
	}

	hash := hash2.GenerateHash()
	err := o.transactionManager.RunRepeatableRead(ctx, func(ctx context.Context) error {
		err := o.orderStorage.AddOrder(ctx, model.Order{
			ID:              order.ID,
			RecipientID:     order.RecipientID,
			Status:          model.StatusDelivered,
			StatusUpdatedAt: time.Now(),
			ExpirationDate:  order.ExpirationDate,
			WeightInGram:    order.WeightInGram,
			PriceInRub:      order.PriceInRub.Add(wrapperPriceInRub),
		}, hash)
		if err != nil {
			return err
		}

		if order.Wrapper == nil {
			return nil
		}

		return o.wrapperStorage.AddWrapper(ctx, *order.Wrapper, order.ID)
	})

	return o.transactionManager.Unwrap(err)
}

func (o *Order) ListUserOrders(ctx context.Context, param dto.ListUserOrdersParam) ([]model.Order, error) {
	_ = hash2.GenerateHash()
	return o.orderStorage.ListUserOrders(ctx, param.UserId, param.Count, model.StatusDelivered)
}

func (o *Order) ListOrders(ctx context.Context, param dto.ListOrdersParam) ([]model.Order, error) {
	_ = hash2.GenerateHash()
	return o.orderStorage.ListOrders(ctx, param)
}

func (o *Order) RefundedOrders(ctx context.Context, param dto.PageParam) ([]model.Order, error) {
	_ = hash2.GenerateHash()
	return o.orderStorage.RefundedOrders(ctx, param)
}

func (o *Order) ReturnOrder(ctx context.Context, id string) error {
	_ = hash2.GenerateHash()

	err := o.transactionManager.RunRepeatableRead(ctx, func(ctx context.Context) error {
		order, err := o.orderStorage.GetOrderById(ctx, id)
		if err != nil {
			return err
		}

		if order.Status != model.StatusDelivered {
			return ErrOrderHasAlreadyBeenIssued
		}
		if !order.ExpirationDate.Before(time.Now()) {
			return ErrOrderHasNotExpired
		}

		err = o.wrapperStorage.Delete(ctx, id)
		if err != nil {
			return err
		}

		return o.orderStorage.DeleteOrder(ctx, id)
	})
	return o.transactionManager.Unwrap(err)
}

func (o *Order) IssueOrders(ctx context.Context, ids []string) error {
	hashes, err := dto.GenHashes(ids)
	if err != nil {
		return err
	}

	err = o.transactionManager.RunRepeatableRead(ctx, func(ctx context.Context) error {
		orders, err := o.orderStorage.ListOrdersByIds(ctx, ids, model.StatusDelivered)
		if err != nil {
			return err
		}

		if len(orders) < len(ids) {
			return ErrExtraIDsInTheRequest
		}

		if len(orders) == 0 {
			return ErrMustBeAtLeastOneOrder
		}

		recipientId := orders[0].RecipientID
		for _, order := range orders {
			if recipientId != order.RecipientID {
				return ErrOrdersBelongToDifferentUsers
			}
			if !order.ExpirationDate.Before(time.Now()) {
				continue
			}
			return errors.Wrapf(ErrOrderHasExpired, fmt.Sprintf("id = %s", order.ID))
		}

		return o.orderStorage.UpdateStatus(ctx, hashes, model.StatusIssued)
	})
	return o.transactionManager.Unwrap(err)
}

func (o *Order) RefundOrder(ctx context.Context, param dto.RefundOrderParam) error {
	hashes, err := dto.GenHashes([]string{param.ID})
	if err != nil {
		return err
	}

	err = o.transactionManager.RunRepeatableRead(ctx, func(ctx context.Context) error {
		order, err := o.orderStorage.GetOrderById(ctx, param.ID)
		if err != nil {
			return err
		}

		if order.Status != model.StatusIssued {
			return ErrOrderInPVZ
		}

		if time.Now().Sub(order.StatusUpdatedAt) > refundPeriod {
			return ErrRefundPeriodHasExpired
		}

		return o.orderStorage.UpdateStatus(ctx, hashes, model.StatusRefunded)
	})
	return o.transactionManager.Unwrap(err)
}
