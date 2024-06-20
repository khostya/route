package model

import (
	"fmt"
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
		ID          string `json:"order_id" db:"id"`
		RecipientID string `json:"recipient_id" db:"recipient_id"`

		Status          Status    `json:"status" db:"status"`
		StatusUpdatedAt time.Time `json:"status_updated_at" db:"status_updated_at"`

		ExpirationDate time.Time `json:"expiration_date" db:"expiration_date"`
	}
)

func (o Order) String() string {
	return fmt.Sprintf("Order(id=%s recipient_id=%s status=%s)", o.ID, o.RecipientID, o.Status)
}
