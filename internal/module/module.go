package module

import (
	"fmt"
	"homework-1/internal/model"
	"homework-1/internal/storage"
	"time"
)

type (
	Storage interface {
		ListOrder(id string, count int, status model.Status) ([]model.Order, error)
		AddOrder(order model.Order) error
		ReWrite(orders []model.Order) error
		ListOrdersByIds(ids []string, delivered model.Status) ([]model.Order, error)
		UpdateStatus(ids []string, issued model.Status) error
		GetOrderById(string) (model.Order, error)
		DeleteOrder(id string) error
		RefundedOrder(get storage.GetParam) ([]model.Order, error)
	}

	Deps struct {
		Storage Storage
	}

	Module struct {
		Deps
	}
)

func NewModule(d Deps) Module {
	return Module{Deps: d}
}

func (m Module) Delivery(order Order) error {
	return m.Storage.AddOrder(model.Order{
		ID:              order.ID,
		RecipientID:     order.RecipientID,
		Status:          model.StatusDelivered,
		StatusUpdatedAt: time.Now(),
		ExpirationDate:  order.ExpirationDate,
	})
}

func (m Module) ListOrder(userID string, count int) ([]model.Order, error) {
	return m.Storage.ListOrder(userID, count, model.StatusDelivered)
}

func (m Module) RefundedOrder(param RefundedOrderParam) ([]model.Order, error) {
	return m.Storage.RefundedOrder(storage.GetParam{Offset: param.Offset, Count: param.Count})
}

func (m Module) ReturnOrder(id string) error {
	order, err := m.Storage.GetOrderById(id)
	if err != nil {
		return err
	}

	if order.Status != model.StatusDelivered {
		return fmt.Errorf("нельзя вернуть заказ, если он был выдан клиенту")
	}
	if !order.ExpirationDate.Before(time.Now()) {
		return fmt.Errorf("у заказа ещё не вышел срок хранения")
	}

	return m.Storage.DeleteOrder(id)
}

func (m Module) IssueOrders(ids []string) error {
	orders, err := m.Storage.ListOrdersByIds(ids, model.StatusDelivered)
	if err != nil {
		return err
	}

	if len(orders) > len(ids) {
		return fmt.Errorf("в запросе присутствуют лишние id")
	}

	for _, order := range orders {
		if !order.ExpirationDate.Before(time.Now()) {
			continue
		}
		return fmt.Errorf("заказ с иднетификатором %s не может быть выдан, "+
			"потому что истек срок хранения", order.ID)
	}

	return m.Storage.UpdateStatus(ids, model.StatusIssued)
}

func (m Module) RefundOrder(param RefundOrderParam) error {
	order, err := m.Storage.GetOrderById(param.ID)
	if err != nil {
		return err
	}

	if order.Status != model.StatusIssued {
		return fmt.Errorf("заказ находится в пвз")
	}

	if order.StatusUpdatedAt.Sub(time.Now()) > (time.Hour * 24 * 2) {
		return fmt.Errorf("заказ не может быть возвращен более чем через два дня")
	}
	return m.Storage.UpdateStatus([]string{param.ID}, model.StatusRefunded)
}
