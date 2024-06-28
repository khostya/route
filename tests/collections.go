//go:build integration

package tests

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

var (
	deliveredOrderWithoutWrapper1 = model.Order{
		ID:              "1",
		RecipientID:     "1",
		Status:          model.StatusDelivered,
		StatusUpdatedAt: time.Now(),
		ExpirationDate:  time.Now().Add(time.Hour * 2),
		WeightInGram:    1,
		PriceInRub:      wrapper.PriceInRub(decimal.NewFromInt(2)),
	}

	deliveredOrder1 = model.Order{
		ID:              "1",
		RecipientID:     "1",
		Status:          model.StatusDelivered,
		StatusUpdatedAt: time.Now(),
		ExpirationDate:  time.Now().Add(time.Hour * 2),
		WeightInGram:    1,
		Wrapper:         wrapper.NewWrapper("box", 1, wrapper.PriceInRub(decimal.NewFromInt(1))),
		PriceInRub:      wrapper.PriceInRub(decimal.NewFromInt(2)),
	}
)
