package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/PetoAdam/homenavi-marketplace/api/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartPostgres(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	container, err := postgres.RunContainer(
		ctx,
		postgres.WithDatabase("marketplace"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithImage("postgres:15-alpine"),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForListeningPort("5432/tcp"),
				wait.ForLog("database system is ready to accept connections"),
			).WithStartupTimeout(90*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("start postgres: %v", err)
	}

	conn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("connection string: %v", err)
	}

	var pool *pgxpool.Pool
	var lastErr error
	for attempt := 0; attempt < 10; attempt++ {
		pool, lastErr = db.Connect(conn)
		if lastErr == nil {
			pingErr := pool.Ping(ctx)
			if pingErr == nil {
				lastErr = nil
				break
			}
			pool.Close()
			lastErr = pingErr
		}
		time.Sleep(500 * time.Millisecond)
	}
	if lastErr != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("db connect: %v", lastErr)
	}

	if err := db.Migrate(ctx, pool); err != nil {
		pool.Close()
		_ = container.Terminate(ctx)
		t.Fatalf("db migrate: %v", err)
	}

	cleanup := func() {
		pool.Close()
		_ = container.Terminate(ctx)
	}

	return pool, cleanup
}
