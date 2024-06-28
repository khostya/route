//go:build integration

package tests

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"homework/internal/model"
	"homework/internal/model/wrapper"
	"homework/internal/storage"
	"homework/internal/storage/transactor"
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
	s.transactor = transactor.NewTransactionManager(db.GetPool())
	s.orderStorage = storage.NewOrderStorage(&s.transactor)
	s.wrapperStorage = storage.NewWrapperStorage(&s.transactor)
	s.ctx = context.Background()
}

func (s *WrapperTestSuite) SetupTest() {
	db.SetUp(s.T(), wrapperTable, orderTable)
}

func (s *WrapperTestSuite) TearDownTest() {
	db.TearDown(s.T())
}

func (s *WrapperTestSuite) TestCreate() {
	err := s.transactor.RunRepeatableRead(s.ctx, func(ctx context.Context) error {
		err := s.orderStorage.AddOrder(ctx, deliveredOrder1, "131")
		if err != nil {
			return err
		}
		return s.wrapperStorage.AddWrapper(ctx, *deliveredOrder1.Wrapper, deliveredOrder1.ID)
	})

	require.Nil(s.T(), err)
}

func (s *WrapperTestSuite) TestGet() {
	err := db.CreateWrapper(s.ctx, deliveredOrder1, "3131")
	require.Nil(s.T(), err)

	wrapper, err := s.get(deliveredOrder1.ID)
	require.Nil(s.T(), err)
	require.Equal(s.T(), *deliveredOrder1.Wrapper, wrapper)
}

func (s *WrapperTestSuite) TestGetWithOrder() {
	err := db.CreateWrapper(s.ctx, deliveredOrder1, "131")
	require.Nil(s.T(), err)

	order, err := s.getOrder(deliveredOrder1)
	require.Nil(s.T(), err)
	require.EqualExportedValues(s.T(), deliveredOrder1, order)
}

func (s *WrapperTestSuite) TestDelete() {
	err := db.CreateWrapper(s.ctx, deliveredOrder1, "131")
	require.Nil(s.T(), err)

	err = s.wrapperStorage.Delete(s.ctx, deliveredOrder1.ID)
	require.Nil(s.T(), err)

	_, err = s.wrapperStorage.GetByOrderId(s.ctx, deliveredOrder1.ID)
	require.ErrorIs(s.T(), storage.ErrNotFound, err)
}

func (s *WrapperTestSuite) get(orderId string) (wrapper.Wrapper, error) {
	return s.wrapperStorage.GetByOrderId(s.ctx, orderId)
}

func (s *WrapperTestSuite) getOrder(order model.Order) (model.Order, error) {
	return s.orderStorage.GetOrderById(s.ctx, order.ID)
}
