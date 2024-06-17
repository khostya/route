package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"homework/internal/model"
	"homework/internal/storage/schema"
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
		UpdateStatus(ctx context.Context, ids schema.IdsWithHashes, issued model.Status) error
		GetOrderById(ctx context.Context, id string) (model.Order, error)
		DeleteOrder(ctx context.Context, id string) error
		RefundedOrders(ctx context.Context, get schema.PageParam) ([]model.Order, error)
	}

	transactionManager interface {
		RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error
		Unwrap(err error) error
	}

	Deps struct {
		Storage            orderStorage
		TransactionManager transactionManager
	}

	Order struct {
		storage            orderStorage
		transactionManager transactionManager
	}
)

func NewOrder(d Deps) Order {
	return Order{storage: d.Storage, transactionManager: d.TransactionManager}
}

func (o *Order) Deliver(ctx context.Context, order DeliverOrderParam) error {
	if order.ExpirationDate.Before(time.Now()) {
		return ErrExpIsNotValid
	}

	hash := hash2.GenerateHash()
	return o.storage.AddOrder(ctx, model.Order{
		ID:              order.ID,
		RecipientID:     order.RecipientID,
		Status:          model.StatusDelivered,
		StatusUpdatedAt: time.Now(),
		ExpirationDate:  order.ExpirationDate,
	}, hash)
}

func (o *Order) ListUserOrders(ctx context.Context, userID string, count uint) ([]model.Order, error) {
	_ = hash2.GenerateHash()
	return o.storage.ListUserOrders(ctx, userID, count, model.StatusDelivered)
}

func (o *Order) RefundedOrders(ctx context.Context, param RefundedOrdersParam) ([]model.Order, error) {
	_ = hash2.GenerateHash()
	return o.storage.RefundedOrders(ctx, schema.PageParam{Page: param.Page, Size: param.Size})
}

func (o *Order) ReturnOrder(ctx context.Context, id string) error {
	_ = hash2.GenerateHash()

	err := o.transactionManager.RunRepeatableRead(ctx, func(ctx context.Context) error {
		order, err := o.storage.GetOrderById(ctx, id)
		if err != nil {
			return err
		}

		if order.Status != model.StatusDelivered {
			return ErrOrderHasAlreadyBeenIssued
		}
		if !order.ExpirationDate.Before(time.Now()) {
			return ErrOrderHasNotExpired
		}

		return o.storage.DeleteOrder(ctx, id)
	})
	return o.transactionManager.Unwrap(err)
}

func (o *Order) IssueOrders(ctx context.Context, ids []string) error {
	hashes, err := o.genHashes(ids)
	if err != nil {
		return err
	}

	err = o.transactionManager.RunRepeatableRead(ctx, func(ctx context.Context) error {
		orders, err := o.storage.ListOrdersByIds(ctx, ids, model.StatusDelivered)
		if err != nil {
			return err
		}

		if len(orders) > len(ids) {
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

		return o.storage.UpdateStatus(ctx, hashes, model.StatusIssued)
	})
	return o.transactionManager.Unwrap(err)
}

func (o *Order) RefundOrder(ctx context.Context, param RefundOrderParam) error {
	hashes, err := o.genHashes([]string{param.ID})
	if err != nil {
		return err
	}

	err = o.transactionManager.RunRepeatableRead(ctx, func(ctx context.Context) error {
		order, err := o.storage.GetOrderById(ctx, param.ID)
		if err != nil {
			return err
		}

		if order.Status != model.StatusIssued {
			return ErrOrderInPVZ
		}

		if order.StatusUpdatedAt.Sub(time.Now()) > refundPeriod {
			return ErrRefundPeriodHasExpired
		}

		return o.storage.UpdateStatus(ctx, hashes, model.StatusRefunded)
	})
	return o.transactionManager.Unwrap(err)
}

func (o *Order) genHashes(strings []string) (schema.IdsWithHashes, error) {
	var hashes []string
	for i := 0; i < len(strings); i++ {
		hashes = append(hashes, hash2.GenerateHash())
	}
	return schema.NewIdsWithHashes(strings, hashes)
}
