//go:build integration

package postgresql

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"homework/internal/cache"
	"homework/internal/dto"
	"homework/internal/model"
	"homework/internal/storage"
	"homework/internal/storage/transactor"
	"homework/tests/postgresql/ids"
	"testing"
)

type OrderTestSuite struct {
	suite.Suite
	ctx          context.Context
	orderStorage *storage.OrderStorage
	transactor   transactor.TransactionManager
}

func TestOrders(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}

func (s *OrderTestSuite) SetupSuite() {
	s.T().Parallel()
	s.transactor = transactor.NewTransactionManager(db.GetPool())
	s.orderStorage = storage.NewOrderStorage(&s.transactor, cache.NewOrdersCache(1, 0))
	s.ctx = context.Background()
}

func (s *OrderTestSuite) SetupTest() {
	s.T().Parallel()
}

func (s *OrderTestSuite) TestCreate() {
	order := NewDeliveredOrderWithoutWrapper(ids.NextID())
	err := s.orderStorage.AddOrder(s.ctx, order, "131")
	require.Nil(s.T(), err)
}

func (s *OrderTestSuite) TestGet() {
	order := NewDeliveredOrderWithoutWrapper(ids.NextID())
	err := db.CreateOrder(s.ctx, order, "131")
	require.Nil(s.T(), err)

	response, err := s.get(order.ID)
	require.EqualExportedValues(s.T(), order, response)
}

func (s *OrderTestSuite) get(id string) (model.Order, error) {
	return s.orderStorage.GetOrderById(s.ctx, id)
}

func (s *OrderTestSuite) TestUpdateStatus() {
	order := NewDeliveredOrderWithoutWrapper(ids.NextID())
	err := db.CreateOrder(s.ctx, order, "131")
	require.Nil(s.T(), err)

	hashes := dto.IdsWithHashes{Ids: []string{order.ID}, Hashes: []string{"311"}}
	err = s.orderStorage.UpdateStatus(s.ctx, hashes, model.StatusIssued)
	require.Nil(s.T(), err)

	response, err := s.get(order.ID)
	require.Nil(s.T(), err)
	require.Equal(s.T(), model.StatusIssued, response.Status)
}

func (s *OrderTestSuite) TestDelete() {
	order := NewDeliveredOrderWithoutWrapper(ids.NextID())
	err := db.CreateOrder(s.ctx, order, "131")
	require.Nil(s.T(), err)

	err = s.orderStorage.DeleteOrder(s.ctx, order.ID)
	require.Nil(s.T(), err)

	_, err = s.get(order.ID)
	require.ErrorIs(s.T(), storage.ErrNotFound, err)
}
