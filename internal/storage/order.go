//go:generate mockgen -source ./mocks/order.go -destination=./mocks/mock_order.go -package=mock_repository
package storage

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"homework/internal/dto"
	"homework/internal/model"
	"homework/internal/storage/schema"
	"homework/internal/storage/transactor"
	"strings"
	"time"
)

const (
	orderTable = "ozon.orders"
	desc       = "DESC"
)

type (
	OrderStorage struct {
		transactor.QueryEngineProvider
		ordersCache ordersCache
	}

	ordersCache interface {
		Get(string) ([]model.Order, bool)
		Put(string, []model.Order)
		RemoveById(string)
		RemoveByIds([]string)
	}
)

func NewOrderStorage(provider transactor.QueryEngineProvider, ordersCache ordersCache) *OrderStorage {
	return &OrderStorage{provider, ordersCache}
}

func (s *OrderStorage) RefundedOrders(ctx context.Context, get dto.PageParam) ([]model.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.OrderStorage.RefundedOrders")
	defer span.Finish()

	offset := get.Size * (get.Page - 1)
	return s.get(ctx, dto.GetParam{Limit: get.Size, Offset: offset, Status: model.StatusRefunded, Order: desc})
}

func (s *OrderStorage) ListOrders(ctx context.Context, get dto.ListOrdersParam) ([]model.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.OrderStorage.ListOrders")
	defer span.Finish()

	offset := get.Size * (get.Page - 1)
	return s.get(ctx, dto.GetParam{
		Limit:       get.Size,
		Offset:      offset,
		Status:      get.Status,
		RecipientId: get.UserId,
		Order:       desc,
	})
}

func (s *OrderStorage) ListUserOrders(ctx context.Context, userId string, count uint, status model.Status) ([]model.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.OrderStorage.ListUserOrders")
	defer span.Finish()

	return s.get(ctx, dto.GetParam{Status: status, Limit: count, RecipientId: userId, Order: desc})
}

func (s *OrderStorage) getByStatus(ctx context.Context, status model.Status) ([]model.Order, error) {
	return s.get(ctx, dto.GetParam{Status: status})
}

func (s *OrderStorage) AddOrder(ctx context.Context, order model.Order, hash string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.OrderStorage.AddOrder")
	defer span.Finish()

	db := s.QueryEngineProvider.GetQueryEngine(ctx)
	record := schema.NewOrder(order, hash)
	query := sq.Insert(orderTable).
		Columns(record.Columns()...).
		Values(record.Values()...).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, rawQuery, args...)
	if err == nil {
		return nil
	}
	if isDuplicateKeyError(err) {
		return ErrDuplicateOrderID
	}
	return err
}

func (s *OrderStorage) ListOrdersByIds(ctx context.Context, ids []string, status model.Status) ([]model.Order, error) {
	return s.get(ctx, dto.GetParam{Ids: ids, Status: status})
}

func (s *OrderStorage) get(ctx context.Context, param dto.GetParam) ([]model.Order, error) {
	cachedOrders, ok := s.ordersCache.Get(param.String())
	if ok {
		return cachedOrders, nil
	}

	db := s.QueryEngineProvider.GetQueryEngine(ctx)
	n := 1

	columns := append(schema.Wrapper{}.SelectColumns(), schema.Order{}.SelectColumns()...)
	query := sq.Select(columns...).
		From(orderTable).
		LeftJoin("ozon.wrappers on wrappers.order_id = orders.id").
		PlaceholderFormat(sq.Dollar)

	if param.Status != "" {
		query = query.Where(fmt.Sprintf("status = $%v", n), param.Status)
		n++
	}
	if param.Ids != nil {
		query = query.Where(fmt.Sprintf("id = ANY($%v)", n), pq.Array(param.Ids))
		n++
	}
	if param.Order != "" {
		query = query.OrderBy(fmt.Sprintf("created_at %v", param.Order))
	}
	if param.RecipientId != "" {
		query = query.Where(fmt.Sprintf("recipient_id = $%v", n), param.RecipientId)
		n++
	}
	if param.Limit != 0 {
		query = query.Limit(uint64(param.Limit))
	}
	if param.Offset != 0 {
		query = query.Offset(uint64(param.Offset))
	}

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return []model.Order{}, err
	}

	var records []schema.WrapperOrder
	if err := pgxscan.Select(ctx, db, &records, rawQuery, args...); err != nil {
		return []model.Order{}, err
	}

	orders, err := schema.ExtractOrdersFromWrapperOrder(records)
	if err != nil {
		return nil, err
	}

	s.ordersCache.Put(param.String(), orders)
	return orders, nil
}

func (s *OrderStorage) UpdateStatus(ctx context.Context, ids dto.IdsWithHashes, status model.Status) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.OrderStorage.UpdateStatus")
	defer span.Finish()

	var setCases strings.Builder
	setCases.WriteString("case\n")
	for i, id := range ids.Ids {
		setCases.WriteString(fmt.Sprintf("when id = '%s' then '%s'\n", id, ids.Hashes[i]))
	}
	setCases.WriteString("end;")

	db := s.QueryEngineProvider.GetQueryEngine(ctx)
	query := sq.Update(orderTable).
		Set("status", status).
		Set("status_updated_at", time.Now()).
		Set("hash", setCases.String()).
		Where("id = ANY($4)", pq.Array(ids.Ids)).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	tag, err := db.Exec(ctx, rawQuery, args...)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	s.ordersCache.RemoveByIds(ids.Ids)
	return nil
}

func (s *OrderStorage) GetOrderById(ctx context.Context, id string) (model.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.OrderStorage.GetOrderById")
	defer span.Finish()

	orders, err := s.get(ctx, dto.GetParam{Ids: []string{id}})
	if err != nil {
		return model.Order{}, err
	}
	if len(orders) != 0 {
		return orders[0], nil
	}
	return model.Order{}, ErrNotFound
}

func (s *OrderStorage) DeleteOrder(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.OrderStorage.DeleteOrder")
	defer span.Finish()

	db := s.QueryEngineProvider.GetQueryEngine(ctx)

	query := sq.Delete(orderTable).
		From(orderTable).
		Where("id = $1", id).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	tag, err := db.Exec(ctx, rawQuery, args...)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	s.ordersCache.RemoveById(id)
	return nil
}
