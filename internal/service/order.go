package service

import (
	"homework/internal/model"
	"homework/internal/storage/schema"
	hash2 "homework/pkg/hash"
	"time"
)

type (
	DeliverOrderParam struct {
		ID          string `json:"order_id"`
		RecipientID string `json:"recipient_id"`

		ExpirationDate time.Time `json:"expiration_date"`
		Wrapper        model.Wrapper
		WeightInKg     float64
	}

	RefundOrderParam struct {
		ID          string `json:"order_id"`
		RecipientID string `json:"recipient_id"`
	}

	RefundedOrdersParam struct {
		Size uint
		Page uint
	}
)

func genHashes(strings []string) (schema.IdsWithHashes, error) {
	var hashes []string
	for i := 0; i < len(strings); i++ {
		hashes = append(hashes, hash2.GenerateHash())
	}
	return schema.NewIdsWithHashes(strings, hashes)
}
