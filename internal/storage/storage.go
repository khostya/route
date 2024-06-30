//go:generate mockgen -source ./storage.go -destination=./mocks/mock_repository.go -package=mock_repository
package storage

import (
	"context"
	"homework/internal/dto"
	"homework/internal/model"
	"homework/internal/model/wrapper"
)

type (
	// orderRepo .
	orderRepo interface {
		AddOrder(ctx context.Context, order model.Order, hash string) error
		RefundedOrders(ctx context.Context, get dto.PageParam) ([]model.Order, error)
		ListUserOrders(ctx context.Context, userId string, count uint, status model.Status) ([]model.Order, error)
		ListOrdersByIds(ctx context.Context, ids []string, status model.Status) ([]model.Order, error)
		UpdateStatus(ctx context.Context, ids dto.IdsWithHashes, status model.Status) error
		GetOrderById(ctx context.Context, id string) (model.Order, error)
		DeleteOrder(ctx context.Context, id string) error
	}

	// wrapperRepo .
	wrapperRepo interface {
		AddWrapper(ctx context.Context, wrapper wrapper.Wrapper, orderId string) error
		Delete(ctx context.Context, orderID string) error
	}
)
