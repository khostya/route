package cli

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	mock_service "homework/internal/service/mocks"
	"testing"
)

type mocks struct {
	mockOrderService *mock_service.MockorderService
}

func newMocks(t *testing.T) mocks {
	ctrl := gomock.NewController(t)

	return mocks{mockOrderService: mock_service.NewMockorderService(ctrl)}
}

func TestExecutor_parceRefundOrder(t *testing.T) {
	t.Parallel()

	type test struct {
		name  string
		input []string
		err   error
	}

	tests := []test{
		{
			name:  ErrIdIsEmpty.Error(),
			input: []string{userIdParamUsage},
			err:   ErrIdIsEmpty,
		},
		{
			name:  ErrUserIsEmpty.Error(),
			input: []string{orderIdParamUsage},
			err:   ErrUserIsEmpty,
		},
		{
			name:  "ok",
			input: []string{orderIdParamUsage, userIdParamUsage},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)

			orderService := newExecutor(mocks.mockOrderService)

			_, err := orderService.parseRefundOrder(tt.input)

			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestExecutor_parseListOrders(t *testing.T) {
	t.Parallel()

	type test struct {
		name  string
		input []string
		err   error
	}

	tests := []test{
		{
			name:  ErrUserIsEmpty.Error(),
			input: []string{sizeParamUsage},
			err:   ErrUserIsEmpty,
		},
		{
			name:  "ok without size",
			input: []string{userIdParamUsage},
		},
		{
			name:  ErrSizeIsNotValid.Error(),
			input: []string{userIdParamUsage, "--size=0"},
			err:   ErrSizeIsNotValid,
		},
		{
			name:  "ok",
			input: []string{userIdParamUsage, sizeParamUsage},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)

			orderService := newExecutor(mocks.mockOrderService)

			_, err := orderService.parseListOrders(tt.input)

			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestExecutor_parseReturnOrder(t *testing.T) {
	t.Parallel()

	type test struct {
		name  string
		input []string
		err   error
	}

	tests := []test{
		{
			name:  ErrIdIsEmpty.Error(),
			input: []string{},
			err:   ErrIdIsEmpty,
		},
		{
			name:  "ok",
			input: []string{orderIdParamUsage},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)

			orderService := newExecutor(mocks.mockOrderService)

			_, err := orderService.parseReturnOrder(tt.input)

			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestExecutor_parseDeliverOrder(t *testing.T) {
	t.Parallel()

	type test struct {
		name  string
		input []string
		err   error
	}

	tests := []test{
		{
			name:  "ok",
			input: []string{userIdParamUsage, orderIdParamUsage, "--wrapper=box", weightInKgUsage, priceInRubParamUsage, expParamUsage},
		},
		{
			name:  ErrExpIsEmpty.Error(),
			input: []string{userIdParamUsage, orderIdParamUsage, "--wrapper=box", weightInKgUsage, priceInRubParamUsage},
			err:   ErrExpIsEmpty,
		},
		{
			name:  ErrUserIsEmpty.Error(),
			input: []string{orderIdParamUsage, "--wrapper=box", weightInKgUsage, priceInRubParamUsage, expParamUsage},
			err:   ErrUserIsEmpty,
		},
		{
			name:  ErrIdIsEmpty.Error(),
			input: []string{userIdParamUsage, "--wrapper=box", weightInKgUsage, priceInRubParamUsage, expParamUsage},
			err:   ErrIdIsEmpty,
		},
		{
			name:  "ok without wrapper",
			input: []string{userIdParamUsage, orderIdParamUsage, weightInKgUsage, priceInRubParamUsage, expParamUsage},
		},
		{
			name:  ErrWeightInKgInNotValid.Error(),
			input: []string{userIdParamUsage, orderIdParamUsage, "--wrapper=box", priceInRubParamUsage, expParamUsage},
			err:   ErrWeightInKgInNotValid,
		},
		{
			name:  ErrPriceInRubIsNotValid.Error(),
			input: []string{userIdParamUsage, orderIdParamUsage, "--wrapper=box", "--price_in_rub=-1", weightInKgUsage, expParamUsage},
			err:   ErrPriceInRubIsNotValid,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)

			orderService := newExecutor(mocks.mockOrderService)

			_, err := orderService.parseDeliverOrder(tt.input)

			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestExecutor_parseListRefunded(t *testing.T) {
	t.Parallel()

	type test struct {
		name  string
		input []string
		err   error
	}

	tests := []test{
		{
			name:  ErrSizeIsNotValid.Error(),
			input: []string{pageParamUsage, "--size=0"},
			err:   ErrSizeIsNotValid,
		},
		{
			name:  ErrPageIsNotValid.Error(),
			input: []string{"--page=0", sizeParamUsage},
			err:   ErrPageIsNotValid,
		},
		{
			name:  "ok without args",
			input: []string{},
		},
		{
			name:  "ok",
			input: []string{sizeParamUsage, pageParamUsage},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)

			orderService := newExecutor(mocks.mockOrderService)

			_, err := orderService.parseListRefunded(tt.input)

			require.ErrorIs(t, err, tt.err)
		})
	}
}
