package transactor

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

const key = "tx"

type (
	QueryEngine interface {
		Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
		Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	}

	QueryEngineProvider interface {
		GetQueryEngine(ctx context.Context) QueryEngine // tx OR pool
	}

	TransactionManager struct {
		pool *pgxpool.Pool
	}
)

func NewTransactionManager(pool *pgxpool.Pool) TransactionManager {
	return TransactionManager{pool}
}

func (tm *TransactionManager) RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error {
	tx, err := tm.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})

	if err != nil {
		return TransactionError{Inner: err}
	}
	if err := fx(context.WithValue(ctx, key, tx)); err != nil {
		return TransactionError{Inner: err, Rollback: tx.Rollback(ctx)}
	}

	if err := tx.Commit(ctx); err != nil {
		return TransactionError{Inner: err, Rollback: tx.Rollback(ctx)}
	}

	return nil
}

func (tm *TransactionManager) Unwrap(err error) error {
	if err == nil {
		return nil
	}

	var transactionError TransactionError
	ok := errors.As(err, &transactionError)
	if !ok {
		return err
	}
	return transactionError.Inner
}

func (tm *TransactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	tx, ok := ctx.Value(key).(QueryEngine)
	if ok && tx != nil {
		return tx
	}

	return tm.pool
}
