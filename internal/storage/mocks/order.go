// DONT EDIT: Auto generated

package mock_repository

import (
	"context"
	"homework/internal/dto"
	"homework/internal/model"
)

// orderStorage ...
type orderStorage interface {
	RefundedOrders(ctx context.Context, get dto.PageParam) ([]model.Order, error)
	ListOrders(ctx context.Context, get dto.ListOrdersParam) ([]model.Order, error)
	ListUserOrders(ctx context.Context, userId string, count uint, status model.Status) ([]model.Order, error)
	AddOrder(ctx context.Context, order model.Order, hash string) error
	ListOrdersByIds(ctx context.Context, ids []string, status model.Status) ([]model.Order, error)
	UpdateStatus(ctx context.Context, ids dto.IdsWithHashes, status model.Status) error
	GetOrderById(ctx context.Context, id string) (model.Order, error)
	DeleteOrder(ctx context.Context, id string) error
}
