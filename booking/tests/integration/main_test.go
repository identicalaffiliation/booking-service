package integration

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	POSTGRES_DOCKER_IMAGE = "postgres:17"
	POSTGRES_DB           = "testdb"
	POSTGRES_USER         = "testuser"
	POSTGRES_PASS         = "testpassword"
	POSTGRES_SSL          = "sslmode=disable"
	MIGRATIONS_DIR        = "./../../../migrator/migrations/booking"
)

var db *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := postgres.Run(
		ctx,
		POSTGRES_DOCKER_IMAGE,
		postgres.WithDatabase(POSTGRES_DB),
		postgres.WithUsername(POSTGRES_USER),
		postgres.WithPassword(POSTGRES_PASS),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("failed to start test postgres container: %v", err)
	}

	connStr, err := container.ConnectionString(ctx, POSTGRES_SSL)
	if err != nil {
		log.Fatalf("failed to get conn str: %v", err)
	}

	db, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to open test postgres pool: %v", err)
	}

	if err := goose.UpContext(ctx, stdlib.OpenDBFromPool(db), MIGRATIONS_DIR); err != nil {
		log.Fatalf("failed to up migrations for test postgres container: %v", err)
	}

	code := m.Run()

	db.Close()
	if err := container.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate test postgres container: %v", err)
	}

	os.Exit(code)
}
