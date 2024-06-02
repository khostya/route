package service

import (
	"fmt"
	"github.com/pkg/errors"
	"homework-1/internal/model"
	"homework-1/internal/storage"
	"time"
)

type (
	orderStorage interface {
		ListOrder(id string, count int, status model.Status) ([]model.Order, error)
		AddOrder(order model.Order) error
		ListOrdersByIds(ids []string, status model.Status) ([]model.Order, error)
		UpdateStatus(ids []string, issued model.Status) error
		GetOrderById(string) (model.Order, error)
		DeleteOrder(id string) error
		RefundedOrder(get storage.GetParam) ([]model.Order, error)
	}

	Deps struct {
		Storage orderStorage
	}

	Order struct {
		storage orderStorage
	}
)

func NewOrder(d Deps) Order {
	return Order{storage: d.Storage}
}

func (m Order) Deliver(order DeliverOrderParam) error {
	return m.storage.AddOrder(model.Order{
		ID:              order.ID,
		RecipientID:     order.RecipientID,
		Status:          model.StatusDelivered,
		StatusUpdatedAt: time.Now(),
		ExpirationDate:  order.ExpirationDate,
	})
}

func (m Order) ListOrder(userID string, count int) ([]model.Order, error) {
	return m.storage.ListOrder(userID, count, model.StatusDelivered)
}

func (m Order) ListRefunded(param RefundedOrderParam) ([]model.Order, error) {
	return m.storage.RefundedOrder(storage.GetParam{Page: param.Page, Size: param.Size})
}

func (m Order) ReturnOrder(id string) error {
	order, err := m.storage.GetOrderById(id)
	if err != nil {
		return err
	}

	if order.Status != model.StatusDelivered {
		return ErrOrderHasAlreadyBeenIssued
	}
	if !order.ExpirationDate.Before(time.Now()) {
		return ErrOrderHasNotExpired
	}

	return m.storage.DeleteOrder(id)
}

func (m Order) IssueOrders(ids []string) error {
	orders, err := m.storage.ListOrdersByIds(ids, model.StatusDelivered)
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

	return m.storage.UpdateStatus(ids, model.StatusIssued)
}

func (m Order) RefundOrder(param RefundOrderParam) error {
	order, err := m.storage.GetOrderById(param.ID)
	if err != nil {
		return err
	}

	if order.Status != model.StatusIssued {
		return ErrOrderInPVZ
	}

	if order.StatusUpdatedAt.Sub(time.Now()) > (time.Hour * 24 * 2) {
		return ErrRefundPeriodHasExpired
	}
	return m.storage.UpdateStatus([]string{param.ID}, model.StatusRefunded)
}
