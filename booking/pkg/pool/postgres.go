package pool

import (
	"context"
	"fmt"

	"github.com/identicalaffiliation/booking-service/booking/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupPool(ctx context.Context, cfg *config.BookingConfig) (*pgxpool.Pool, error, func()) {
	pool, err := pgxpool.New(ctx, cfg.DbUrl)
	if err != nil {
		return nil, fmt.Errorf("open new pgx pool: %w", err), func() {}
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping pgx pool: %w", err), func() {}
	}

	pool.Config().MaxConns = cfg.MaxConns
	pool.Config().MaxConnLifetime = cfg.MaxLifetime

	return pool, nil, func() {
		pool.Close()
	}
}
