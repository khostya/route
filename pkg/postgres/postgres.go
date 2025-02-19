package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"

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

func PoolFromEnv(ctx context.Context, key string) (*pgxpool.Pool, error) {
	url := os.Getenv(key)
	if url == "" {
		return nil, errors.New(fmt.Sprintf("Unable to parse %s", "DATABASE_URL"))
	}

	return Pool(ctx, url)
}
