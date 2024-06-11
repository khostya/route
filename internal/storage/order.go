package storage

import (
	"homework/internal/model"
	"homework/pkg/hash"
)

type (
	record struct {
		model.Order
		Hash string `json:"hash"`
	}

	GetParam struct {
		Size int
		Page int
	}
)

func newRecord(orders []model.Order) []record {
	return mapFunc(orders, func(order model.Order) record {
		return record{order, hash.GenerateHash()}
	})
}

func extractOrders(records []record) []model.Order {
	return mapFunc(records, func(record record) model.Order {
		return record.Order
	})
}

func mapFunc[IN any, OUT any](in []IN, m func(IN) OUT) []OUT {
	var out []OUT

	for _, i := range in {
		out = append(out, m(i))
	}

	return out
}
