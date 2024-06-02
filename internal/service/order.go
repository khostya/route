package service

import "time"

type (
	DeliverOrderParam struct {
		ID          string `json:"order_id"`
		RecipientID string `json:"recipient_id"`

		ExpirationDate time.Time `json:"expiration_date"`
	}

	RefundOrderParam struct {
		ID          string `json:"order_id"`
		RecipientID string `json:"recipient_id"`
	}

	RefundedOrderParam struct {
		Size int
		Page int
	}
)
