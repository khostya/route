package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	pool "homework/pkg/postgres"
	"os"
	"strings"
	"testing"
)

type DBPool struct {
	pool *pgxpool.Pool
}

func NewFromEnv() *DBPool {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		panic("TEST_DATABASE_URL isn`t set")
	}

	pool, err := pool.Pool(context.Background(), url)
	if err != nil {
		panic(err)
	}
	return &DBPool{pool: pool}
}

func (d *DBPool) SetUp(t *testing.T, tableName ...string) {
	t.Helper()
	d.truncateTable(context.Background(), tableName...)
}

func (d *DBPool) TearDown(t *testing.T) {
	t.Helper()
}

func (d *DBPool) truncateTable(ctx context.Context, tableName ...string) {

	q := fmt.Sprintf("TRUNCATE %s", strings.Join(tableName, ","))
	if _, err := d.pool.Exec(ctx, q); err != nil {
		panic(err)
	}
}

func (d *DBPool) GetPool() *pgxpool.Pool {
	return d.pool
}

func (d *DBPool) Close() {
	d.pool.Close()
}
