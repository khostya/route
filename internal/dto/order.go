package dto

import (
	"fmt"
	"homework/internal/model"
	"homework/internal/model/wrapper"
	"strings"
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

func (p ListOrdersParam) String() string {
	return fmt.Sprintf("[ListOrdersParam]: page=%v userID=%v size=%v status=%v ", p.Page, p.UserId, p.Size, string(p.Status))
}

func (p GetParam) String() string {
	return fmt.Sprintf("[GetParam]: ids=%v status=%v order=%v limit=%v RecipientId=%v Offset=%v",
		strings.Join(p.Ids, ", "), string(p.Status), p.Order, p.Limit, p.RecipientId, p.Offset)
}
