//go:generate mockgen -source ./order.go -destination=./mocks/mock_order.go -package=mock_service
package service

import (
	"context"
	"homework/internal/model"
	"homework/internal/model/wrapper"
	"homework/internal/storage/schema"
	hash2 "homework/pkg/hash"
	"time"
)

type (
	orderService interface {
		Deliver(ctx context.Context, order DeliverOrderParam) error
		ListUserOrders(ctx context.Context, param ListUserOrdersParam) ([]model.Order, error)
		RefundedOrders(ctx context.Context, param RefundedOrdersParam) ([]model.Order, error)
		ReturnOrder(ctx context.Context, id string) error
		IssueOrders(ctx context.Context, ids []string) error
		RefundOrder(ctx context.Context, param RefundOrderParam) error
	}

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
