package storage

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"homework/internal/model"
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

func (w *WrapperStorage) AddWrapper(ctx context.Context, wrapper model.Wrapper, orderId string) error {
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
