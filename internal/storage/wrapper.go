//go:generate mockgen -source ./mocks/wrapper.go -destination=./mocks/mock_wrapper.go -package=mock_repository
package storage

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/opentracing/opentracing-go"
	"homework/internal/model/wrapper"
	"homework/internal/storage/schema"
	"homework/internal/storage/transactor"
)

const (
	wrapperTable = "ozon.wrappers"
)

type (
	WrapperStorage struct {
		transactor.QueryEngineProvider
	}
)

func NewWrapperStorage(provider transactor.QueryEngineProvider) *WrapperStorage {
	return &WrapperStorage{provider}
}

func (w *WrapperStorage) AddWrapper(ctx context.Context, wrapper wrapper.Wrapper, orderId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.WrapperStorage.AddWrapper")
	defer span.Finish()

	db := w.QueryEngineProvider.GetQueryEngine(ctx)
	record := schema.NewWrapper(wrapper, orderId)
	query := sq.Insert(wrapperTable).
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

func (w *WrapperStorage) Delete(ctx context.Context, orderId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.WrapperStorage.Delete")
	defer span.Finish()

	db := w.QueryEngineProvider.GetQueryEngine(ctx)

	query := sq.Delete(wrapperTable).
		From(wrapperTable).
		Where("order_id = $1", orderId).
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

func (w *WrapperStorage) GetByOrderId(ctx context.Context, orderId string) (wrapper.Wrapper, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.WrapperStorage.GetByOrderId")
	defer span.Finish()

	db := w.QueryEngineProvider.GetQueryEngine(ctx)

	columns := append(schema.Wrapper{}.SelectColumns())
	query := sq.Select(columns...).
		From(orderTable).
		LeftJoin("ozon.wrappers on wrappers.order_id = orders.id").
		Where("order_id = $1", orderId).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return wrapper.Wrapper{}, err
	}

	var records []schema.Wrapper
	if err := pgxscan.Select(ctx, db, &records, rawQuery, args...); err != nil {
		return wrapper.Wrapper{}, err
	}
	if len(records) == 0 {
		return wrapper.Wrapper{}, ErrNotFound
	}

	return schema.ExtractWrapper(records[0]), nil
}
