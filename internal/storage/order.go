package storage

import (
	"homework-1/internal/model"
	"homework-1/pkg/hash"
)

type (
	Record struct {
		Hash   string
		Orders []model.Order
	}

	GetParam struct {
		Count  int
		Offset int
	}
)

func newRecord(orders []model.Order) Record {
	return Record{
		Hash:   hash.GenerateHash(),
		Orders: orders,
	}
}
