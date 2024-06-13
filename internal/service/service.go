package service

import (
	"fmt"
	"github.com/pkg/errors"
	"homework/internal/model"
	"homework/internal/storage"
	"sync"
	"time"
)

const (
	RefundPeriod = time.Hour * 24 * 2
)

type (
	orderStorage interface {
		ListUserOrders(id string, count int, status model.Status) ([]model.Order, error)
		AddOrder(order model.Order) error
		ListOrdersByIds(ids []string, status model.Status) ([]model.Order, error)
		UpdateStatus(ids []string, issued model.Status) error
		GetOrderById(string) (model.Order, error)
		DeleteOrder(id string) error
		RefundedOrders(get storage.GetParam) ([]model.Order, error)
	}

	Deps struct {
		Storage orderStorage
	}

	Order struct {
		storage orderStorage
		mutex   sync.RWMutex
	}
)

func NewOrder(d Deps) Order {
	return Order{storage: d.Storage}
}

func (o *Order) Deliver(order DeliverOrderParam) error {
	if order.ExpirationDate.Before(time.Now()) {
		return ErrExpIsNotValid
	}

	o.mutex.Lock()
	defer o.mutex.Unlock()

	return o.storage.AddOrder(model.Order{
		ID:              order.ID,
		RecipientID:     order.RecipientID,
		Status:          model.StatusDelivered,
		StatusUpdatedAt: time.Now(),
		ExpirationDate:  order.ExpirationDate,
	})
}

func (o *Order) ListUserOrders(userID string, count int) ([]model.Order, error) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	return o.storage.ListUserOrders(userID, count, model.StatusDelivered)
}

func (o *Order) RefundedOrders(param RefundedOrdersParam) ([]model.Order, error) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	return o.storage.RefundedOrders(storage.GetParam{Page: param.Page, Size: param.Size})
}

func (o *Order) ReturnOrder(id string) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	order, err := o.storage.GetOrderById(id)
	if err != nil {
		return err
	}

	if order.Status != model.StatusDelivered {
		return ErrOrderHasAlreadyBeenIssued
	}
	if !order.ExpirationDate.Before(time.Now()) {
		return ErrOrderHasNotExpired
	}

	return o.storage.DeleteOrder(id)
}

func (o *Order) IssueOrders(ids []string) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	orders, err := o.storage.ListOrdersByIds(ids, model.StatusDelivered)
	if err != nil {
		return err
	}

	if len(orders) > len(ids) {
		return ErrExtraIDsInTheRequest
	}

	for _, order := range orders {
		if !order.ExpirationDate.Before(time.Now()) {
			continue
		}
		return errors.Wrapf(ErrOrderHasExpired, fmt.Sprintf("id = %s", order.ID))
	}

	return o.storage.UpdateStatus(ids, model.StatusIssued)
}

func (o *Order) RefundOrder(param RefundOrderParam) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	order, err := o.storage.GetOrderById(param.ID)
	if err != nil {
		return err
	}

	if order.Status != model.StatusIssued {
		return ErrOrderInPVZ
	}

	if order.StatusUpdatedAt.Sub(time.Now()) > RefundPeriod {
		return ErrRefundPeriodHasExpired
	}

	return o.storage.UpdateStatus([]string{param.ID}, model.StatusRefunded)
}
