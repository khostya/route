package storage

import (
	"encoding/json"
	"errors"
	"homework/internal/model"
	"os"
	"slices"
	"sync"
	"time"
)

type (
	Storage struct {
		fileName string
		mutex    sync.RWMutex
	}
)

func NewStorage(fileName string) (*Storage, error) {
	storage := Storage{fileName: fileName}
	if _, err := os.Stat(fileName); !errors.Is(err, os.ErrNotExist) {
		return &Storage{fileName: fileName}, nil
	}

	if errCreateFile := storage.createFile(); errCreateFile != nil {
		return nil, errCreateFile
	}

	err := storage.reWrite([]record{})
	if err != nil {
		return nil, err
	}

	return &Storage{fileName: fileName}, nil
}

func (s *Storage) RefundedOrders(get GetParam) ([]model.Order, error) {
	s.mutex.RLock()
	orders, err := s.getByStatus(model.StatusRefunded)
	s.mutex.RUnlock()

	if err != nil {
		return nil, err
	}

	left := min(get.Page*get.Size, len(orders))

	right := min(get.Page*get.Size+get.Size, len(orders))
	if right < 0 {
		right = len(orders)
	}

	return orders[left:right], nil
}

func (s *Storage) ListUserOrders(userId string, count int, status model.Status) ([]model.Order, error) {
	s.mutex.RLock()
	orders, err := s.getByStatus(status)
	s.mutex.RUnlock()

	if err != nil {
		return nil, err
	}

	orders = slices.DeleteFunc(orders, func(order model.Order) bool {
		return order.RecipientID != userId
	})

	return orders[max(len(orders)-count, 0):], nil
}

func (s *Storage) getByStatus(status model.Status) ([]model.Order, error) {
	records, err := s.allRecords()
	if err != nil {
		return nil, err
	}

	records = slices.DeleteFunc(records, func(order record) bool {
		return order.Status != status
	})
	return extractOrders(records), nil
}

func (s *Storage) allRecords() ([]record, error) {
	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return nil, err
	}

	var records []record
	err = json.Unmarshal(b, &records)

	return records, err
}

func (s *Storage) AddOrder(order model.Order, hash string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	records, err := s.allRecords()
	if err != nil {
		return err
	}

	isDuplicate := slices.ContainsFunc(records, func(o record) bool {
		return order.ID == o.ID
	})
	if isDuplicate {
		return ErrDuplicateOrderID
	}

	records = append(records, newRecord(order, hash))
	return s.reWrite(records)
}

func (s *Storage) reWrite(records []record) error {
	bWrite, errMarshal := json.MarshalIndent(records, "  ", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	return os.WriteFile(s.fileName, bWrite, 0666)
}

func (s *Storage) ListOrdersByIds(ids []string, status model.Status) ([]model.Order, error) {
	s.mutex.RLock()
	orders, err := s.getByStatus(status)
	s.mutex.RUnlock()

	if err != nil {
		return nil, err
	}

	setIds := toSet(ids)

	orders = slices.DeleteFunc(orders, func(order model.Order) bool {
		return !setIds[order.ID]
	})

	return orders, nil
}

func (s *Storage) UpdateStatus(ids ListWithHashes, status model.Status) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	orders, err := s.allRecords()
	if err != nil {
		return err
	}

	setIds := toSet(ids.list)
	for i := range orders {
		if !setIds[orders[i].ID] {
			continue
		}
		orders[i].Status = status
		orders[i].StatusUpdatedAt = time.Now()
		orders[i].Hash = ids.hashes[i]
	}

	return s.reWrite(orders)
}

func (s *Storage) GetOrderById(id string) (model.Order, error) {
	s.mutex.RLock()
	records, err := s.allRecords()
	s.mutex.RUnlock()

	if err != nil {
		return model.Order{}, err
	}

	for _, record := range records {
		if record.ID == id {
			return record.Order, nil
		}
	}
	return model.Order{}, ErrNotFound
}

func (s *Storage) DeleteOrder(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	records, err := s.allRecords()
	if err != nil {
		return err
	}

	records = slices.DeleteFunc(records, func(record record) bool {
		return record.ID == id
	})

	return s.reWrite(records)
}

func (s *Storage) createFile() error {
	f, err := os.Create(s.fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func toSet[T comparable](s []T) map[T]bool {
	var m = make(map[T]bool, len(s))
	for _, el := range s {
		m[el] = true
	}
	return m
}
