package model

import (
	"fmt"
	"homework/internal/model/wrapper"
	"time"
)

var (
	StatusDelivered = Status("delivered")
	StatusIssued    = Status("issued")
	StatusRefunded  = Status("refunded")
	TimeFormat      = time.RFC3339
)

type (
	Status string

	Order struct {
		ID          string `json:"order_id"`
		RecipientID string `json:"recipient_id"`

		Status          Status    `json:"status"`
		StatusUpdatedAt time.Time `json:"status_updated_at"`

		ExpirationDate time.Time `json:"expiration_date"`
		WeightInGram   float64   `json:"weight_in_gram"`
		Wrapper        *wrapper.Wrapper
		PriceInRub     wrapper.PriceInRub
	}
)

func (o Order) String() string {
	return fmt.Sprintf("Order(id=%s recipient_id=%s status=%s)", o.ID, o.RecipientID, o.Status)
}
