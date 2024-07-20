//go:build integration

package postgresql

import (
	"github.com/shopspring/decimal"
	"homework/internal/model"
	"homework/internal/model/wrapper"
	"time"
)

const (
	orderTable   = "ozon.orders"
	wrapperTable = "ozon.wrappers"
)

func NewDeliveredOrderWithoutWrapper(id string) model.Order {
	return model.Order{
		ID:              id,
		RecipientID:     "1",
		Status:          model.StatusDelivered,
		StatusUpdatedAt: time.Now(),
		ExpirationDate:  time.Now().Add(time.Hour * 2),
		WeightInGram:    1,
		PriceInRub:      wrapper.PriceInRub(decimal.NewFromInt(2)),
	}
}

func NewDeliveredOrder(id string) model.Order {
	return model.Order{
		ID:              id,
		RecipientID:     "1",
		Status:          model.StatusDelivered,
		StatusUpdatedAt: time.Now(),
		ExpirationDate:  time.Now().Add(time.Hour * 2),
		WeightInGram:    1,
		Wrapper:         wrapper.NewWrapper("box", 1, wrapper.PriceInRub(decimal.NewFromInt(1))),
		PriceInRub:      wrapper.PriceInRub(decimal.NewFromInt(2)),
	}
}
