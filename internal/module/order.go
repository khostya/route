package module

import "time"

type (
	Order struct {
		ID          string `json:"order_id"`
		RecipientID string `json:"recipient_id"`

		ExpirationDate time.Time `json:"expiration_date"`
	}

	RefundOrderParam struct {
		ID          string `json:"order_id"`
		RecipientID string `json:"recipient_id"`
	}

	RefundedOrderParam struct {
		Count  int
		Offset int
	}
)
