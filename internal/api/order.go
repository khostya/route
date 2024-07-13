package api

import (
	"context"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework/internal/dto"
	"homework/internal/model"
	"homework/internal/model/wrapper"
	"homework/pkg/api/order/v1"
)

type (
	OrderService struct {
		service orderService
		order.UnimplementedOrderServer
	}

	orderService interface {
		Deliver(ctx context.Context, order dto.DeliverOrderParam) error
		ListOrders(ctx context.Context, param dto.ListOrdersParam) ([]model.Order, error)
		ReturnOrder(ctx context.Context, id string) error
		IssueOrders(ctx context.Context, ids []string) error
		RefundOrder(ctx context.Context, param dto.RefundOrderParam) error
	}
)

func NewOrderService(orderService orderService) *OrderService {
	return &OrderService{service: orderService}
}

func (o *OrderService) ReturnOrder(ctx context.Context, req *order.ReturnOrderRequest) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := o.service.ReturnOrder(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (o *OrderService) IssueOrders(ctx context.Context, req *order.IssueOrdersRequest) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := o.service.IssueOrders(ctx, req.GetIds())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (o *OrderService) RefundOrder(ctx context.Context, req *order.RefundOrderRequest) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := o.service.RefundOrder(ctx, dto.RefundOrderParam{
		ID:          req.GetOrderID(),
		RecipientID: req.GetUserID(),
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (o *OrderService) ListOrders(ctx context.Context, req *order.ListOrdersRequest) (*order.ListOrdersResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	orders, err := o.service.ListOrders(ctx, dto.ListOrdersParam{
		UserId: req.GetUserID(),
		Size:   uint(req.GetSize()),
		Page:   uint(req.GetPage()),
		Status: model.Status(order.OrderStatus_name[int32(req.GetStatus())]),
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var resp order.ListOrdersResponse
	for _, o := range orders {
		statusValue, ok := order.OrderStatus_value[string(o.Status)]
		if !ok {
			return nil, status.Error(codes.Internal, err.Error())
		}

		respOrder := &order.ListOrdersResponse_Order{
			RecipientID: o.RecipientID,
			Id:          o.ID,
			Status:      order.OrderStatus(statusValue),
		}
		resp.Orders = append(resp.Orders, respOrder)
	}

	return &resp, nil
}

func (o *OrderService) DeliverOrder(ctx context.Context, req *order.DeliverOrderRequest) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	priceInRub := wrapper.PriceInRub(decimal.NewFromFloat(float64(req.GetPriceInRub())))
	wrapper, err := wrapper.NewDefaultWrapper(wrapper.WrapperType(req.GetWrapperType()))
	if req.WrapperType != nil && err != nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, err.Error())
	}

	err = o.service.Deliver(ctx, dto.DeliverOrderParam{
		ID:             req.GetOrderID(),
		RecipientID:    req.GetUserID(),
		ExpirationDate: req.GetExp().AsTime(),
		WeightInGram:   float64(req.GetWeightInKg() * 1000),
		Wrapper:        wrapper,
		PriceInRub:     priceInRub,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
