package service

import (
	"fmt"
	"github.com/pkg/errors"
	"homework/internal/model"
	"homework/internal/storage"
	hash2 "homework/pkg/hash"
	"sync"
	"time"
)

const (
	RefundPeriod = time.Hour * 24 * 2
)

type (
	orderStorage interface {
		ListUserOrders(id string, count int, status model.Status) ([]model.Order, error)
		AddOrder(order model.Order, hash string) error
		ListOrdersByIds(ids []string, status model.Status) ([]model.Order, error)
		UpdateStatus(ids storage.ListWithHashes, issued model.Status) error
		GetOrderById(id string) (model.Order, error)
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

	hash := hash2.GenerateHash()

	o.mutex.Lock()
	defer o.mutex.Unlock()

	return o.storage.AddOrder(model.Order{
		ID:              order.ID,
		RecipientID:     order.RecipientID,
		Status:          model.StatusDelivered,
		StatusUpdatedAt: time.Now(),
		ExpirationDate:  order.ExpirationDate,
	}, hash)
}

func (o *Order) ListUserOrders(userID string, count int) ([]model.Order, error) {
	_ = hash2.GenerateHash()

	o.mutex.RLock()
	defer o.mutex.RUnlock()

	return o.storage.ListUserOrders(userID, count, model.StatusDelivered)
}

func (o *Order) RefundedOrders(param RefundedOrdersParam) ([]model.Order, error) {
	_ = hash2.GenerateHash()

	o.mutex.RLock()
	defer o.mutex.RUnlock()

	return o.storage.RefundedOrders(storage.GetParam{Page: param.Page, Size: param.Size})
}

func (o *Order) ReturnOrder(id string) error {
	_ = hash2.GenerateHash()

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
	hashes, err := o.genHashes(ids)
	if err != nil {
		return err
	}

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

	return o.storage.UpdateStatus(hashes, model.StatusIssued)
}

func (o *Order) RefundOrder(param RefundOrderParam) error {
	hashes, err := o.genHashes([]string{param.ID})
	if err != nil {
		return err
	}

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

	return o.storage.UpdateStatus(hashes, model.StatusRefunded)
}

func (o *Order) genHashes(strings []string) (storage.ListWithHashes, error) {
	var hashes []string
	for i := 0; i < len(strings); i++ {
		hashes = append(hashes, hash2.GenerateHash())
	}
	return storage.NewListWithHashes(strings, hashes)
}
