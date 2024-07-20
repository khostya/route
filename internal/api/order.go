package api

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework/internal/dto"
	"homework/internal/metrics"
	"homework/internal/model"
	"homework/internal/model/wrapper"
	"homework/internal/service"
	"homework/internal/storage"
	"homework/pkg/api/order/v1"
)

type (
	OrderService struct {
		service orderService
		cache   ordersCache
		order.UnimplementedOrderServer
	}

	orderService interface {
		Deliver(ctx context.Context, order dto.DeliverOrderParam) error
		ListOrders(ctx context.Context, param dto.ListOrdersParam) ([]model.Order, error)
		ReturnOrder(ctx context.Context, id string) error
		IssueOrders(ctx context.Context, ids []string) error
		RefundOrder(ctx context.Context, param dto.RefundOrderParam) error
	}

	ordersCache interface {
		Get(string) ([]model.Order, bool)
		Put(string, []model.Order)
		RemoveById(string) bool
	}
)

func NewOrderService(orderService orderService, cache ordersCache) *OrderService {
	return &OrderService{
		service: orderService,
		cache:   cache,
	}
}

func (o *OrderService) ReturnOrder(ctx context.Context, req *order.ReturnOrderRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.OrderService.ReturnOrder")
	defer span.Finish()

	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := o.service.ReturnOrder(ctx, req.GetId())

	if err := toGRPCError(err); err != nil {
		return nil, err
	}

	o.cache.RemoveById(req.GetId())
	return &emptypb.Empty{}, nil
}

func (o *OrderService) IssueOrders(ctx context.Context, req *order.IssueOrdersRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.OrderService.IssueOrders")
	defer span.Finish()

	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := o.service.IssueOrders(ctx, req.GetIds())
	if err := toGRPCError(err); err != nil {
		return nil, err
	}

	for _, id := range req.Ids {
		o.cache.RemoveById(id)
	}
	metrics.AddIssuedOrders(len(req.Ids))
	return &emptypb.Empty{}, nil
}

func (o *OrderService) RefundOrder(ctx context.Context, req *order.RefundOrderRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.OrderService.RefundOrder")
	defer span.Finish()

	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := o.service.RefundOrder(ctx, dto.RefundOrderParam{
		ID:          req.GetOrderID(),
		RecipientID: req.GetUserID(),
	})

	if err := toGRPCError(err); err != nil {
		return nil, err
	}

	o.cache.RemoveById(req.GetOrderID())
	return &emptypb.Empty{}, nil
}

func (o *OrderService) ListOrders(ctx context.Context, req *order.ListOrdersRequest) (*order.ListOrdersResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.OrderService.ListOrders")
	defer span.Finish()

	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	param := dto.ListOrdersParam{
		UserId: req.GetUserID(),
		Size:   uint(req.GetSize()),
		Page:   uint(req.GetPage()),
		Status: grpcOrderStatusToDomain(req.GetStatus()),
	}

	cachedOrders, ok := o.cache.Get(param.String())
	if ok {
		return o.buildListOrderResp(cachedOrders), nil
	}

	orders, err := o.service.ListOrders(ctx, param)
	if err := toGRPCError(err); err != nil {
		return nil, err
	}

	o.cache.Put(param.String(), orders)
	resp := o.buildListOrderResp(orders)
	return resp, nil
}

func (o *OrderService) buildListOrderResp(orders []model.Order) *order.ListOrdersResponse {
	var resp order.ListOrdersResponse
	for _, o := range orders {
		respOrder := &order.ListOrdersResponse_Order{
			RecipientID: o.RecipientID,
			Id:          o.ID,
			Status:      domainOrderStatusToGRPC(o.Status),
		}
		resp.Orders = append(resp.Orders, respOrder)
	}
	return &resp
}

func (o *OrderService) DeliverOrder(ctx context.Context, req *order.DeliverOrderRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.OrderService.DeliverOrder")
	defer span.Finish()

	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	priceInRub := wrapper.PriceInRub(decimal.NewFromFloat(float64(req.GetPriceInRub())))
	wrapper, err := wrapper.NewDefaultWrapper(grpcWrapperTypeToDomain(req.GetWrapperType()))
	if req.WrapperType != order.WrapperType_WRAPPER_TYPE_NONE && err != nil {
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

	if err := toGRPCError(err); err != nil {
		return nil, err
	}

	o.cache.RemoveById(req.GetOrderID())
	return &emptypb.Empty{}, nil
}

func toGRPCError(err error) error {
	var orderServiceError service.OrderServiceError
	isOrderServiceError := errors.As(err, &orderServiceError)
	if isOrderServiceError {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	switch {
	case errors.Is(err, storage.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, storage.ErrDuplicateOrderID):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		if err == nil {
			return nil
		}
		return status.Error(codes.Internal, err.Error())
	}
}

func grpcOrderStatusToDomain(orderStatus order.OrderStatus) model.Status {
	return map[order.OrderStatus]model.Status{
		order.OrderStatus_ORDER_STATUS_ANY:       model.StatusNone,
		order.OrderStatus_ORDER_STATUS_DELIVERED: model.StatusDelivered,
		order.OrderStatus_ORDER_STATUS_ISSUED:    model.StatusIssued,
		order.OrderStatus_ORDER_STATUS_REFUNDED:  model.StatusRefunded,
	}[orderStatus]
}

func domainOrderStatusToGRPC(orderStatus model.Status) order.OrderStatus {
	return map[model.Status]order.OrderStatus{
		model.StatusDelivered: order.OrderStatus_ORDER_STATUS_DELIVERED,
		model.StatusIssued:    order.OrderStatus_ORDER_STATUS_ISSUED,
		model.StatusRefunded:  order.OrderStatus_ORDER_STATUS_REFUNDED,
	}[orderStatus]
}

func grpcWrapperTypeToDomain(wrapperType order.WrapperType) wrapper.WrapperType {
	return map[order.WrapperType]wrapper.WrapperType{
		order.WrapperType_WRAPPER_TYPE_NONE:    wrapper.NoneWrapper,
		order.WrapperType_WRAPPER_TYPE_BOX:     wrapper.BoxWrapper,
		order.WrapperType_WRAPPER_TYPE_PACKAGE: wrapper.PackageWrapper,
		order.WrapperType_WRAPPER_TYPE_STRETCH: wrapper.StretchWrapper,
	}[wrapperType]
}
