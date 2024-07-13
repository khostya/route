package api

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"homework/internal/model"
	mock_service "homework/internal/service/mocks"
	"homework/pkg/api/order/v1"
	"testing"
	"time"
)

type mocks struct {
	mockOrderService *mock_service.MockorderService
}

func newMocks(t *testing.T) mocks {
	ctrl := gomock.NewController(t)

	return mocks{
		mockOrderService: mock_service.NewMockorderService(ctrl),
	}
}

func TestDeliver(t *testing.T) {
	t.Parallel()
	var (
		randomWrapper = "31331fdsf"
		boxWrapper    = "box"
	)

	type test struct {
		name   string
		input  *order.DeliverOrderRequest
		code   codes.Code
		mockFn func(m mocks)
	}
	var ctx = context.Background()
	tests := []test{
		{
			name: "validate error",
			input: &order.DeliverOrderRequest{
				OrderID:    "1",
				UserID:     "1",
				Exp:        timestamppb.New(time.Now().Add(-time.Hour)),
				PriceInRub: 1.0,
				WeightInKg: 1.0,
			},
			code: codes.InvalidArgument,
			mockFn: func(m mocks) {
			},
		},
		{
			name: "wrapper does not exists",
			input: &order.DeliverOrderRequest{
				OrderID:     "1",
				UserID:      "1",
				WrapperType: &randomWrapper,
				Exp:         timestamppb.New(time.Now().Add(-time.Hour)),
				PriceInRub:  1.0,
				WeightInKg:  1.0,
			},
			code: codes.InvalidArgument,
			mockFn: func(m mocks) {
			},
		},
		{
			name: "ok",
			input: &order.DeliverOrderRequest{
				OrderID:     "1",
				UserID:      "1",
				WrapperType: &boxWrapper,
				Exp:         timestamppb.New(time.Now().Add(time.Hour)),
				PriceInRub:  1.0,
				WeightInKg:  1.0,
			},
			code: codes.OK,
			mockFn: func(m mocks) {
				m.mockOrderService.EXPECT().Deliver(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)
			tt.mockFn(mocks)

			service := NewOrderService(mocks.mockOrderService)
			_, err := service.DeliverOrder(ctx, tt.input)
			status, ok := status.FromError(err)
			if ok && tt.code == codes.OK {
				return
			}
			require.Equal(t, tt.code, status.Code())
		})
	}
}

func TestListOrders(t *testing.T) {
	t.Parallel()
	var (
		userID        = "1"
		size   uint32 = 1
		order1        = order.ListOrdersResponse_Order{RecipientID: "1", Id: "1", Status: order.OrderStatus_Delivered}
	)

	type test struct {
		name    string
		input   *order.ListOrdersRequest
		code    codes.Code
		mockFn  func(m mocks)
		result  *order.ListOrdersResponse
		wantErr bool
	}
	var ctx = context.Background()
	tests := []test{
		{
			name: "ok",
			input: &order.ListOrdersRequest{
				UserID: &userID,
				Size:   &size,
			},
			code: codes.OK,
			mockFn: func(m mocks) {
				orders := []model.Order{
					{
						Status:      model.Status(order.OrderStatus_name[int32(order1.Status)]),
						ID:          order1.Id,
						RecipientID: order1.RecipientID,
					},
				}
				m.mockOrderService.EXPECT().ListOrders(gomock.Any(), gomock.Any()).Times(1).Return(orders, nil)
			},
			wantErr: false,
			result:  &order.ListOrdersResponse{Orders: []*order.ListOrdersResponse_Order{&order1}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)
			tt.mockFn(mocks)

			service := NewOrderService(mocks.mockOrderService)
			orders, err := service.ListOrders(ctx, tt.input)
			status, _ := status.FromError(err)

			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, orders.Orders, []*order.ListOrdersResponse_Order{&order1})
			require.Equal(t, tt.code, status.Code())
		})
	}
}
