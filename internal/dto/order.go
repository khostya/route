package dto

import (
	"homework/internal/model"
	"homework/internal/model/wrapper"
	"time"
)

type (
	DeliverOrderParam struct {
		ID          string `json:"order_id"`
		RecipientID string `json:"recipient_id"`

		ExpirationDate time.Time `json:"expiration_date"`
		Wrapper        *wrapper.Wrapper
		WeightInGram   float64
		PriceInRub     wrapper.PriceInRub
	}

	RefundOrderParam struct {
		ID          string `json:"order_id"`
		RecipientID string `json:"recipient_id"`
	}

	ListUserOrdersParam struct {
		UserId string
		Count  uint
	}

	ListOrdersParam struct {
		UserId string
		Size   uint
		Page   uint
		Status model.Status
	}

	PageParam struct {
		Size uint
		Page uint
	}

	GetParam struct {
		Ids         []string
		Status      model.Status
		Order       string
		Limit       uint
		RecipientId string
		Offset      uint
	}
)
