package sqldb

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/goosemigrator"
	"gluttony/migrations"
	"testing"
)

func NewTestPGXPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	m := goosemigrator.New(".", goosemigrator.WithFS(migrations.FS))
	cfg := pgtestdb.Custom(t, pgtestdb.Config{
		DriverName: "pgx",
		Host:       "localhost",
		Port:       "5432",
		User:       "dev",
		Password:   "dev",
		Database:   "dev",
		Options:    "sslmode=disable",
	}, m)

	pool, err := ConnectPostgres(context.Background(), Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Database,
		Options:  cfg.Options,
	})
	if err != nil {
		t.Fatalf("sql.ConnectPostgres: %v", err)
	}

	t.Cleanup(pool.Close)

	return pool
}
