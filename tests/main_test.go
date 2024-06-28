//go:build integration

package tests

import (
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

	db.Close()
	os.Exit(code)
}
