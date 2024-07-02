package service

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"homework/internal/dto"
	"homework/internal/model"
	"homework/internal/model/wrapper"
	mock_repository "homework/internal/storage/mocks"
	mock_transactor "homework/internal/storage/transactor/mocks"
	"testing"
	"time"
)

type mocks struct {
	mockOrderRepository   *mock_repository.MockorderRepo
	mockWrapperRepository *mock_repository.MockwrapperRepo
	mockTransactor        *mock_transactor.MockTransactor
}

func newMocks(t *testing.T) mocks {
	ctrl := gomock.NewController(t)

	return mocks{
		mockTransactor:        mock_transactor.NewMockTransactor(ctrl),
		mockWrapperRepository: mock_repository.NewMockwrapperRepo(ctrl),
		mockOrderRepository:   mock_repository.NewMockorderRepo(ctrl),
	}
}

func TestOrderService_Deliver(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  dto.DeliverOrderParam
		err    error
		mockFn func(m mocks)
	}
	var ctx = context.Background()
	tests := []test{
		{
			name: "exp is not valid",
			input: dto.DeliverOrderParam{
				ID:             "1",
				RecipientID:    "1",
				ExpirationDate: time.Now().Add(-time.Minute * 10),
				Wrapper:        nil,
				WeightInGram:   0,
				PriceInRub:     wrapper.PriceInRub(decimal.NewFromInt(0)),
			},
			err: ErrExpIsNotValid,
			mockFn: func(m mocks) {
			},
		},
		{
			name: "order weight greater than wrapper capacity",
			input: dto.DeliverOrderParam{
				ID:             "1",
				RecipientID:    "1",
				ExpirationDate: time.Now().Add(time.Minute * 10),
				Wrapper:        wrapper.NewWrapper("box", 1, wrapper.PriceInRub(decimal.NewFromFloat(0))),
				WeightInGram:   2,
				PriceInRub:     wrapper.PriceInRub(decimal.NewFromInt(0)),
			},
			err: ErrOrderWeightGreaterThanWrapperCapacity,
			mockFn: func(m mocks) {

			},
		},
		{
			name: "ok",
			input: dto.DeliverOrderParam{
				ID:             "1",
				RecipientID:    "1",
				ExpirationDate: time.Now().Add(time.Minute * 10),
				Wrapper:        wrapper.NewWrapper("box", 20, wrapper.PriceInRub(decimal.NewFromFloat(0))),
				WeightInGram:   10,
				PriceInRub:     wrapper.PriceInRub(decimal.NewFromInt(0)),
			},
			mockFn: func(m mocks) {
				m.mockOrderRepository.EXPECT().AddOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
				m.mockWrapperRepository.EXPECT().AddWrapper(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
				m.mockTransactor.EXPECT().Unwrap(nil).Times(1).Return(nil)
				m.mockTransactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
						return transaction(ctx)
					})
			},
		},
		{
			name: "ok without wrapper",
			input: dto.DeliverOrderParam{
				ID:             "1",
				RecipientID:    "1",
				ExpirationDate: time.Now().Add(time.Minute * 10),
				WeightInGram:   2,
				PriceInRub:     wrapper.PriceInRub(decimal.NewFromInt(0)),
			},
			mockFn: func(m mocks) {
				m.mockOrderRepository.EXPECT().AddOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
				m.mockTransactor.EXPECT().Unwrap(nil).Times(1).Return(nil)
				m.mockTransactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
						return transaction(ctx)
					})
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)
			tt.mockFn(mocks)
			orderService := NewOrder(Deps{
				WrapperStorage:     mocks.mockWrapperRepository,
				Storage:            mocks.mockOrderRepository,
				TransactionManager: mocks.mockTransactor,
			})

			err := orderService.Deliver(ctx, tt.input)

			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestOrderService_ReturnOrder(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  string
		err    error
		mockFn func(m mocks)
	}

	var ctx = context.Background()
	tests := []test{
		{
			name:  "has already been issued",
			input: "1",
			mockFn: func(m mocks) {
				m.mockOrderRepository.EXPECT().GetOrderById(gomock.Any(), gomock.Any()).
					Times(1).Return(model.Order{Status: model.StatusIssued}, nil)

				m.mockTransactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
						err := transaction(ctx)
						m.mockTransactor.EXPECT().Unwrap(gomock.Any()).Times(1).Return(err)
						return err
					})
			},
			err: ErrOrderHasAlreadyBeenIssued,
		},
		{
			name:  "order has not expired",
			input: "1",
			mockFn: func(m mocks) {
				m.mockOrderRepository.EXPECT().GetOrderById(gomock.Any(), gomock.Any()).
					Times(1).Return(model.Order{
					Status:         model.StatusDelivered,
					ExpirationDate: time.Now().Add(time.Hour),
				}, nil)

				m.mockTransactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
						err := transaction(ctx)
						m.mockTransactor.EXPECT().Unwrap(gomock.Any()).Times(1).Return(err)
						return err
					})
			},
			err: ErrOrderHasNotExpired,
		},
		{
			name:  "ok",
			input: "1",
			mockFn: func(m mocks) {
				m.mockOrderRepository.EXPECT().GetOrderById(gomock.Any(), gomock.Any()).
					Times(1).Return(model.Order{
					Status:         model.StatusDelivered,
					ExpirationDate: time.Now().Add(-time.Hour),
				}, nil)

				m.mockOrderRepository.EXPECT().DeleteOrder(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				m.mockWrapperRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				m.mockTransactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
						err := transaction(ctx)
						m.mockTransactor.EXPECT().Unwrap(gomock.Any()).Times(1).Return(err)
						return err
					})
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)
			tt.mockFn(mocks)
			orderService := NewOrder(Deps{
				WrapperStorage:     mocks.mockWrapperRepository,
				Storage:            mocks.mockOrderRepository,
				TransactionManager: mocks.mockTransactor,
			})

			err := orderService.ReturnOrder(ctx, tt.input)

			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestOrderService_RefundOrder(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  dto.RefundOrderParam
		err    error
		mockFn func(m mocks)
	}

	var ctx = context.Background()
	tests := []test{
		{
			name:  "order in pvz",
			input: dto.RefundOrderParam{ID: "1", RecipientID: "1"},
			mockFn: func(m mocks) {
				m.mockOrderRepository.EXPECT().GetOrderById(gomock.Any(), gomock.Any()).
					Times(1).Return(model.Order{Status: model.StatusDelivered}, nil)
				m.mockTransactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
						err := transaction(ctx)
						m.mockTransactor.EXPECT().Unwrap(gomock.Any()).Times(1).Return(err)
						return err
					})
			},
			err: ErrOrderInPVZ,
		},
		{
			name:  "refund period has expired",
			input: dto.RefundOrderParam{ID: "1", RecipientID: "1"},
			mockFn: func(m mocks) {
				m.mockOrderRepository.EXPECT().GetOrderById(gomock.Any(), gomock.Any()).
					Times(1).Return(model.Order{
					Status:          model.StatusIssued,
					StatusUpdatedAt: time.Now().Add(-2 * refundPeriod),
				}, nil)
				m.mockTransactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
						err := transaction(ctx)
						m.mockTransactor.EXPECT().Unwrap(gomock.Any()).Times(1).Return(err)
						return err
					})
			},
			err: ErrRefundPeriodHasExpired,
		},
		{
			name:  "ok",
			input: dto.RefundOrderParam{ID: "1", RecipientID: "1"},
			mockFn: func(m mocks) {
				m.mockOrderRepository.EXPECT().GetOrderById(gomock.Any(), gomock.Any()).
					Times(1).Return(model.Order{
					Status:          model.StatusIssued,
					StatusUpdatedAt: time.Now(),
				}, nil)
				m.mockOrderRepository.EXPECT().UpdateStatus(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return(nil)
				m.mockTransactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
						return transaction(ctx)
					})
				m.mockTransactor.EXPECT().Unwrap(nil).Times(1).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)
			tt.mockFn(mocks)
			orderService := NewOrder(Deps{
				WrapperStorage:     mocks.mockWrapperRepository,
				Storage:            mocks.mockOrderRepository,
				TransactionManager: mocks.mockTransactor,
			})

			err := orderService.RefundOrder(ctx, tt.input)

			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestOrderService_IssueOrders(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  []string
		err    error
		mockFn func(m mocks)
	}

	var ctx = context.Background()
	tests := []test{
		{
			name:  "extra IDs in the request",
			input: []string{"1", "2"},
			mockFn: func(m mocks) {
				m.mockOrderRepository.EXPECT().ListOrdersByIds(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return([]model.Order{
					{Status: model.StatusDelivered},
				}, nil)

				m.mockTransactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
						err := transaction(ctx)
						m.mockTransactor.EXPECT().Unwrap(gomock.Any()).Times(1).Return(err)
						return err
					})
			},
			err: ErrExtraIDsInTheRequest,
		},
		{
			name:  "must be at least one order",
			input: []string{},
			mockFn: func(m mocks) {
				m.mockOrderRepository.EXPECT().ListOrdersByIds(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return([]model.Order{}, nil)

				m.mockTransactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
						err := transaction(ctx)
						m.mockTransactor.EXPECT().Unwrap(gomock.Any()).Times(1).Return(err)
						return err
					})
			},
			err: ErrMustBeAtLeastOneOrder,
		},
		{
			name:  "orders belong to different users",
			input: []string{"1", "2"},
			mockFn: func(m mocks) {
				m.mockOrderRepository.EXPECT().ListOrdersByIds(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return([]model.Order{
					{Status: model.StatusDelivered, RecipientID: "1", ID: "1", ExpirationDate: time.Now().Add(time.Hour)},
					{Status: model.StatusDelivered, RecipientID: "2", ID: "2", ExpirationDate: time.Now().Add(time.Hour)},
				}, nil)

				m.mockTransactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
						err := transaction(ctx)
						m.mockTransactor.EXPECT().Unwrap(gomock.Any()).Times(1).Return(err)
						return err
					})
			},
			err: ErrOrdersBelongToDifferentUsers,
		},
		{
			name:  "order has expired",
			input: []string{"1", "2"},
			mockFn: func(m mocks) {
				m.mockOrderRepository.EXPECT().ListOrdersByIds(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return([]model.Order{
					{Status: model.StatusDelivered, RecipientID: "1", ID: "1", ExpirationDate: time.Now().Add(-time.Hour)},
					{Status: model.StatusDelivered, RecipientID: "1", ID: "2", ExpirationDate: time.Now().Add(time.Hour)},
				}, nil)

				m.mockTransactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
						err := transaction(ctx)
						m.mockTransactor.EXPECT().Unwrap(gomock.Any()).Times(1).Return(err)
						return err
					})
			},
			err: ErrOrderHasExpired,
		},
		{
			name:  "ok",
			input: []string{"1", "2"},
			mockFn: func(m mocks) {
				m.mockOrderRepository.EXPECT().ListOrdersByIds(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return([]model.Order{
					{Status: model.StatusDelivered, RecipientID: "1", ID: "1", ExpirationDate: time.Now().Add(time.Hour)},
					{Status: model.StatusDelivered, RecipientID: "1", ID: "2", ExpirationDate: time.Now().Add(time.Hour)},
				}, nil)

				m.mockOrderRepository.EXPECT().UpdateStatus(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)

				m.mockTransactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
						err := transaction(ctx)
						m.mockTransactor.EXPECT().Unwrap(gomock.Any()).Times(1).Return(err)
						return err
					})
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)
			tt.mockFn(mocks)
			orderService := NewOrder(Deps{
				WrapperStorage:     mocks.mockWrapperRepository,
				Storage:            mocks.mockOrderRepository,
				TransactionManager: mocks.mockTransactor,
			})

			err := orderService.IssueOrders(ctx, tt.input)

			require.ErrorIs(t, err, tt.err)
		})
	}
}
