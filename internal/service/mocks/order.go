// DONT EDIT: Auto generated

package mock_service

import (
	"context"
	"homework/internal/dto"
	"homework/internal/model"
)

// orderService ...
type orderService interface {
	Deliver(ctx context.Context, order dto.DeliverOrderParam) error
	ListUserOrders(ctx context.Context, param dto.ListUserOrdersParam) ([]model.Order, error)
	ListOrders(ctx context.Context, param dto.ListOrdersParam) ([]model.Order, error)
	RefundedOrders(ctx context.Context, param dto.PageParam) ([]model.Order, error)
	ReturnOrder(ctx context.Context, id string) error
	IssueOrders(ctx context.Context, ids []string) error
	RefundOrder(ctx context.Context, param dto.RefundOrderParam) error
}
