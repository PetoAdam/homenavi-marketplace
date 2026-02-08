package db

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := fmt.Sprintf("migrations/%s", entry.Name())
		sqlBytes, err := migrationFiles.ReadFile(path)
		if err != nil {
			return err
		}
		if _, err := pool.Exec(ctx, string(sqlBytes)); err != nil {
			return fmt.Errorf("migration %s failed: %w", entry.Name(), err)
		}
	}
	return nil
}
