//go:build integration

package postgresql

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"homework/internal/model"
	"homework/internal/model/wrapper"
	"homework/internal/storage"
	"homework/internal/storage/transactor"
	"homework/tests/postgresql/ids"
	"testing"
)

type WrapperTestSuite struct {
	suite.Suite
	ctx            context.Context
	orderStorage   *storage.OrderStorage
	wrapperStorage *storage.WrapperStorage
	transactor     transactor.TransactionManager
}

func TestWrapper(t *testing.T) {
	suite.Run(t, new(WrapperTestSuite))
}

func (s *WrapperTestSuite) SetupSuite() {
	s.T().Parallel()
	s.transactor = transactor.NewTransactionManager(db.GetPool())
	s.orderStorage = storage.NewOrderStorage(&s.transactor)
	s.wrapperStorage = storage.NewWrapperStorage(&s.transactor)
	s.ctx = context.Background()
}

func (s *WrapperTestSuite) SetupTest() {
	s.T().Parallel()
}

func (s *WrapperTestSuite) TestCreate() {
	order := NewDeliveredOrder(ids.NextID())
	err := s.transactor.RunRepeatableRead(s.ctx, func(ctx context.Context) error {
		err := s.orderStorage.AddOrder(ctx, order, "131")
		if err != nil {
			return err
		}
		return s.wrapperStorage.AddWrapper(ctx, *order.Wrapper, order.ID)
	})

	require.Nil(s.T(), err)
}

func (s *WrapperTestSuite) TestGet() {
	order := NewDeliveredOrder(ids.NextID())
	err := db.CreateWrapper(s.ctx, order, "3131")
	require.Nil(s.T(), err)

	wrapper, err := s.get(order.ID)
	require.Nil(s.T(), err)
	require.Equal(s.T(), *order.Wrapper, wrapper)
}

func (s *WrapperTestSuite) TestGetWithOrder() {
	order := NewDeliveredOrder(ids.NextID())
	err := db.CreateWrapper(s.ctx, order, "131")
	require.Nil(s.T(), err)

	response, err := s.getOrder(order)
	require.Nil(s.T(), err)
	require.EqualExportedValues(s.T(), order, response)
}

func (s *WrapperTestSuite) TestDelete() {
	order := NewDeliveredOrder(ids.NextID())
	err := db.CreateWrapper(s.ctx, order, "131")
	require.Nil(s.T(), err)

	err = s.wrapperStorage.Delete(s.ctx, order.ID)
	require.Nil(s.T(), err)

	_, err = s.wrapperStorage.GetByOrderId(s.ctx, order.ID)
	require.ErrorIs(s.T(), storage.ErrNotFound, err)
}

func (s *WrapperTestSuite) get(orderId string) (wrapper.Wrapper, error) {
	return s.wrapperStorage.GetByOrderId(s.ctx, orderId)
}

func (s *WrapperTestSuite) getOrder(order model.Order) (model.Order, error) {
	return s.orderStorage.GetOrderById(s.ctx, order.ID)
}
