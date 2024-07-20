//go:generate mockgen -source ./mocks/order.go -destination=./mocks/mock_order.go -package=mock_service
package service

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
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

	OrderService struct {
		orderStorage       orderStorage
		transactionManager transactionManager
		wrapperStorage     wrapperStorage
	}
)

func NewOrder(d Deps) OrderService {
	return OrderService{
		orderStorage:       d.Storage,
		transactionManager: d.TransactionManager,
		wrapperStorage:     d.WrapperStorage,
	}
}

func (o *OrderService) Deliver(ctx context.Context, param dto.DeliverOrderParam) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.OrderService.Deliver")
	defer span.Finish()

	hash := hash2.GenerateHash()

	return o.deliver(ctx, param, hash)
}

func (o *OrderService) deliver(ctx context.Context, param dto.DeliverOrderParam, hash string) error {
	if param.ExpirationDate.Before(time.Now()) {
		return ErrExpIsNotValid
	}
	if param.Wrapper != nil && !param.Wrapper.WillFitGram(param.WeightInGram) {
		message := fmt.Sprintf("capacity_in_gram = %v", param.Wrapper.GetCapacityInGram())
		return errors.Wrap(ErrOrderWeightGreaterThanWrapperCapacity, message)
	}

	wrapperPriceInRub := wrapper.PriceInRub(decimal.NewFromInt(0))
	if param.Wrapper != nil {
		wrapperPriceInRub = param.Wrapper.GetPriceInRub()
	}

	err := o.transactionManager.RunRepeatableRead(ctx, func(ctx context.Context) error {
		err := o.orderStorage.AddOrder(ctx, model.Order{
			ID:              param.ID,
			RecipientID:     param.RecipientID,
			Status:          model.StatusDelivered,
			StatusUpdatedAt: time.Now(),
			ExpirationDate:  param.ExpirationDate,
			WeightInGram:    param.WeightInGram,
			PriceInRub:      param.PriceInRub.Add(wrapperPriceInRub),
		}, hash)
		if err != nil {
			return err
		}

		if param.Wrapper == nil {
			return nil
		}

		return o.wrapperStorage.AddWrapper(ctx, *param.Wrapper, param.ID)
	})

	return o.transactionManager.Unwrap(err)
}

func (o *OrderService) ListUserOrders(ctx context.Context, param dto.ListUserOrdersParam) ([]model.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.OrderService.ListUserOrders")
	defer span.Finish()

	_ = hash2.GenerateHash()
	return o.listUserOrders(ctx, param)
}

func (o *OrderService) listUserOrders(ctx context.Context, param dto.ListUserOrdersParam) ([]model.Order, error) {
	return o.orderStorage.ListUserOrders(ctx, param.UserId, param.Count, model.StatusDelivered)
}

func (o *OrderService) ListOrders(ctx context.Context, param dto.ListOrdersParam) ([]model.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.OrderService.ListOrders")
	defer span.Finish()

	_ = hash2.GenerateHash()
	return o.listOrders(ctx, param)
}

func (o *OrderService) listOrders(ctx context.Context, param dto.ListOrdersParam) ([]model.Order, error) {
	return o.orderStorage.ListOrders(ctx, param)
}

func (o *OrderService) RefundedOrders(ctx context.Context, param dto.PageParam) ([]model.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.OrderService.RefundedOrders")
	defer span.Finish()

	_ = hash2.GenerateHash()
	return o.refundedOrders(ctx, param)
}

func (o *OrderService) refundedOrders(ctx context.Context, param dto.PageParam) ([]model.Order, error) {
	return o.orderStorage.RefundedOrders(ctx, param)
}

func (o *OrderService) ReturnOrder(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.OrderService.ReturnOrder")
	defer span.Finish()

	_ = hash2.GenerateHash()

	return o.returnOrder(ctx, id)
}

func (o *OrderService) returnOrder(ctx context.Context, id string) error {
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

func (o *OrderService) IssueOrders(ctx context.Context, ids []string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.OrderService.IssueOrders")
	defer span.Finish()

	hashes, err := dto.GenHashes(ids)
	if err != nil {
		return err
	}

	return o.issueOrders(ctx, ids, hashes)
}

func (o *OrderService) issueOrders(ctx context.Context, ids []string, hashes dto.IdsWithHashes) error {
	err := o.transactionManager.RunRepeatableRead(ctx, func(ctx context.Context) error {
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

func (o *OrderService) RefundOrder(ctx context.Context, param dto.RefundOrderParam) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.OrderService.RefundOrder")
	defer span.Finish()

	hashes, err := dto.GenHashes([]string{param.ID})
	if err != nil {
		return err
	}

	return o.refundOrder(ctx, param, hashes)
}

func (o *OrderService) refundOrder(ctx context.Context, param dto.RefundOrderParam, hashes dto.IdsWithHashes) error {
	err := o.transactionManager.RunRepeatableRead(ctx, func(ctx context.Context) error {
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
