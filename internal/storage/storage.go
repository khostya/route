package storage

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/lib/pq"
	"homework/internal/model"
	"homework/internal/storage/schema"
	"homework/internal/storage/transactor"
	"time"
)

const (
	orderTable = "ozon.orders"
)

var (
	desc = "DESC"
)

type (
	Storage struct {
		transactor.QueryEngineProvider
	}
)

func NewStorage(provider transactor.QueryEngineProvider) *Storage {
	return &Storage{provider}
}

func (s *Storage) RefundedOrders(ctx context.Context, get schema.PageParam) ([]model.Order, error) {
	offset := get.Size * get.Page
	return s.get(ctx, schema.GetParam{Limit: &get.Size, Offset: &offset, Status: &model.StatusRefunded, Order: &desc})
}

func (s *Storage) ListUserOrders(ctx context.Context, userId string, count int, status model.Status) ([]model.Order, error) {
	return s.get(ctx, schema.GetParam{Status: &status, Limit: &count, RecipientId: &userId, Order: &desc})
}

func (s *Storage) getByStatus(ctx context.Context, status model.Status) ([]model.Order, error) {
	return s.get(ctx, schema.GetParam{Status: &status})
}

func (s *Storage) AddOrder(ctx context.Context, order model.Order, hash string) error {
	db := s.QueryEngineProvider.GetQueryEngine(ctx)
	record := schema.NewRecord(order, hash)
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

func (s *Storage) ListOrdersByIds(ctx context.Context, ids []string, status model.Status) ([]model.Order, error) {
	return s.get(ctx, schema.GetParam{Ids: ids, Status: &status})
}

func (s *Storage) get(ctx context.Context, param schema.GetParam) ([]model.Order, error) {
	db := s.QueryEngineProvider.GetQueryEngine(ctx)
	n := 1

	query := sq.Select(schema.Record{}.Columns()...).
		From(orderTable).
		PlaceholderFormat(sq.Dollar)
	if param.Status != nil {
		query = query.Where(fmt.Sprintf("status = $%v", n), *param.Status)
		n++
	}
	if param.Ids != nil {
		query = query.Where(fmt.Sprintf("id = ANY($%v)", n), pq.Array(param.Ids))
		n++
	}
	if param.Order != nil {
		query = query.OrderBy(fmt.Sprintf("created_at %v", *param.Order))
	}
	if param.RecipientId != nil {
		query = query.Where(fmt.Sprintf("recipient_id = $%v", n), *param.RecipientId)
		n++
	}
	if param.Limit != nil {
		query = query.Limit(uint64(*param.Limit))
	}
	if param.Offset != nil {
		query = query.Offset(uint64(*param.Offset))
	}

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return []model.Order{}, err
	}

	var records []schema.Record
	if err := pgxscan.Select(ctx, db, &records, rawQuery, args...); err != nil {
		return []model.Order{}, err
	}

	return schema.ExtractOrders(records), nil
}

func (s *Storage) UpdateStatus(ctx context.Context, ids []string, status model.Status, hash string) error {
	db := s.QueryEngineProvider.GetQueryEngine(ctx)
	query := sq.Update(orderTable).
		Set("status", status).
		Set("status_updated_at", time.Now()).
		Set("hash", hash).
		Where("id = ANY($4)", pq.Array(ids)).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	tag, err := db.Exec(ctx, rawQuery, args...)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

func (s *Storage) GetOrderById(ctx context.Context, id string) (model.Order, error) {
	orders, err := s.get(ctx, schema.GetParam{Ids: []string{id}})
	if err != nil {
		return model.Order{}, err
	}
	if len(orders) != 0 {
		return orders[0], nil
	}
	return model.Order{}, ErrNotFound
}

func (s *Storage) DeleteOrder(ctx context.Context, id string) error {
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
	return err
}
