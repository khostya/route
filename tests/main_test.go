//go:build integration

package tests

import (
	"context"
	"homework/tests/postgresql"
	"os"
	"testing"
)

var (
	db *postgresql.DBPool
)

func TestMain(m *testing.M) {
	db = postgresql.NewFromEnv()

	code := m.Run()

	db.TruncateTable(context.Background(), wrapperTable, orderTable)
	db.Close()

	os.Exit(code)
}
