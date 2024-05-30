package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"homework-1/internal/model"
	"os"
	"slices"
	"time"
)

type (
	Storage struct {
		fileName string
	}

	Order struct {
		ID          string `json:"order_id"`
		RecipientID string `json:"recipient_id"`

		ExpirationDate time.Time `json:"expiration_date"`
		Status         model.Status
	}
)

func NewStorage(fileName string) (Storage, error) {
	storage := Storage{fileName: fileName}
	if _, err := os.Stat(fileName); !errors.Is(err, os.ErrNotExist) {
		return Storage{fileName: fileName}, nil
	}

	if errCreateFile := storage.createFile(); errCreateFile != nil {
		return storage, errCreateFile
	}

	err := storage.ReWrite([]model.Order{})
	if err != nil {
		return storage, err
	}

	return Storage{fileName: fileName}, nil
}

func (s Storage) RefundedOrder(get GetParam) ([]model.Order, error) {
	orders, err := s.getByStatus(model.StatusRefunded)
	if err != nil {
		return nil, err
	}

	right := min(get.Offset+get.Count, len(orders))
	if right < 0 {
		right = len(orders)
	}

	left := min(get.Offset, len(orders))

	return orders[left:right], nil
}

func (s Storage) ListOrder(userId string, count int, status model.Status) ([]model.Order, error) {
	orders, err := s.getByStatus(status)
	if err != nil {
		return nil, err
	}

	orders = deleteAll(orders, func(order model.Order) bool {
		return order.RecipientID != userId
	})

	return orders[max(len(orders)-count, 0):], nil
}

func (s Storage) getByStatus(status model.Status) ([]model.Order, error) {
	orders, err := s.allOrders()
	if err != nil {
		return nil, err
	}

	orders = deleteAll(orders, func(order model.Order) bool {
		return order.Status != status
	})
	return orders, nil
}

func (s Storage) allOrders() ([]model.Order, error) {
	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return nil, err
	}

	var record Record
	err = json.Unmarshal(b, &record)

	return record.Orders, err
}

func (s Storage) AddOrder(order model.Order) error {
	orders, err := s.allOrders()
	if err != nil {
		return err
	}

	isDuplicate := slices.ContainsFunc(orders, func(o model.Order) bool {
		return order.ID == o.ID
	})
	if isDuplicate {
		return fmt.Errorf("duplicate order id")
	}

	orders = append(orders, order)
	return s.ReWrite(orders)
}

func (s Storage) ReWrite(orders []model.Order) error {
	record := newRecord(orders)
	bWrite, errMarshal := json.MarshalIndent(record, "  ", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	return os.WriteFile(s.fileName, bWrite, 0666)
}

func (s Storage) ListOrdersByIds(ids []string, status model.Status) ([]model.Order, error) {
	orders, err := s.getByStatus(status)
	if err != nil {
		return nil, err
	}

	setIds := toCounter(ids)

	orders = deleteAll(orders, func(order model.Order) bool {
		return !setIds[order.ID]
	})

	return orders, nil
}

func (s Storage) UpdateStatus(ids []string, status model.Status) error {
	orders, err := s.allOrders()
	if err != nil {
		return err
	}

	setIds := toCounter(ids)
	for i := range orders {
		if !setIds[orders[i].ID] {
			continue
		}
		orders[i].Status = status
		orders[i].StatusUpdatedAt = time.Now()
	}

	return s.ReWrite(orders)
}

func (s Storage) GetOrderById(id string) (model.Order, error) {
	orders, err := s.allOrders()
	if err != nil {
		return model.Order{}, err
	}

	for _, order := range orders {
		if order.ID == id {
			return order, nil
		}
	}
	return model.Order{}, fmt.Errorf("not found")
}

func (s Storage) DeleteOrder(id string) error {
	orders, err := s.allOrders()
	if err != nil {
		return err
	}

	orders = slices.DeleteFunc(orders, func(order model.Order) bool {
		return order.ID == id
	})

	return s.ReWrite(orders)
}

func (s Storage) createFile() error {
	f, err := os.Create(s.fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func toCounter[T comparable](s []T) map[T]bool {
	var m = make(map[T]bool, len(s))
	for _, el := range s {
		m[el] = true
	}
	return m
}

func deleteAll[S ~[]E, E any](s S, del func(E) bool) S {
	result := make(S, 0, len(s))

	for _, v := range s {
		if del(v) {
			continue
		}
		result = append(result, v)
	}

	return result
}
