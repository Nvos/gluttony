package testdb

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
)

func Seed(t *testing.T, pool *pgxpool.Pool, name string) {
	t.Helper()

	script, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatalf("read sql script file: %v", err)
	}

	_, err = pool.Exec(context.Background(), string(script))
	if err != nil {
		t.Fatalf("exec sql script file: %v", err)
	}
}
