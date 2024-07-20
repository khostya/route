package postgres

import (
	"context"
	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

func Pool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create connection pool")
	}

	pool.Config().ConnConfig.LogLevel = pgx.LogLevelDebug
	return pool, nil
}
