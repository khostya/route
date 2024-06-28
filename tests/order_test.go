//go:build integration

package tests

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"homework/internal/model"
	"homework/internal/storage"
	"homework/internal/storage/schema"
	"homework/internal/storage/transactor"
	"testing"
)

type OrderTestSuite struct {
	suite.Suite
	ctx          context.Context
	orderStorage *storage.OrderStorage
}

func TestOrders(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}

func (s *OrderTestSuite) SetupSuite() {
	transactor := transactor.NewTransactionManager(db.GetPool())
	s.orderStorage = storage.NewOrderStorage(&transactor)
	s.ctx = context.Background()
}

func (s *OrderTestSuite) SetupTest() {
	db.SetUp(s.T(), wrapperTable, orderTable)
}

func (s *OrderTestSuite) TearDownTest() {
	db.TearDown(s.T())
}

func (s *OrderTestSuite) TestCreate() {
	err := s.orderStorage.AddOrder(s.ctx, deliveredOrderWithoutWrapper1, "131")
	require.Nil(s.T(), err)
}

func (s *OrderTestSuite) TestGet() {
	err := db.CreateOrder(s.ctx, deliveredOrderWithoutWrapper1, "131")
	require.Nil(s.T(), err)

	order, err := s.get(deliveredOrderWithoutWrapper1.ID)
	require.EqualExportedValues(s.T(), deliveredOrderWithoutWrapper1, order)
}

func (s *OrderTestSuite) get(id string) (model.Order, error) {
	return s.orderStorage.GetOrderById(s.ctx, id)
}

func (s *OrderTestSuite) TestUpdateStatus() {
	err := db.CreateOrder(s.ctx, deliveredOrderWithoutWrapper1, "131")
	require.Nil(s.T(), err)

	hashes := schema.IdsWithHashes{Ids: []string{deliveredOrderWithoutWrapper1.ID}, Hashes: []string{"311"}}
	err = s.orderStorage.UpdateStatus(s.ctx, hashes, model.StatusIssued)
	require.Nil(s.T(), err)

	order, err := s.get(deliveredOrderWithoutWrapper1.ID)
	require.Nil(s.T(), err)
	require.Equal(s.T(), model.StatusIssued, order.Status)
}

func (s *OrderTestSuite) TestDelete() {
	err := db.CreateOrder(s.ctx, deliveredOrderWithoutWrapper1, "131")
	require.Nil(s.T(), err)

	err = s.orderStorage.DeleteOrder(s.ctx, deliveredOrderWithoutWrapper1.ID)
	require.Nil(s.T(), err)

	_, err = s.get(deliveredOrderWithoutWrapper1.ID)
	require.ErrorIs(s.T(), storage.ErrNotFound, err)
}
