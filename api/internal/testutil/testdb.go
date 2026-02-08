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
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp").WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		t.Fatalf("start postgres: %v", err)
	}

	conn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("connection string: %v", err)
	}

	pool, err := db.Connect(conn)
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("db connect: %v", err)
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
