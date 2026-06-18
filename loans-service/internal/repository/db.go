package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, ErrDatabaseURLNotSet
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConnectionPool, err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDatabaseUnreachable, err)
	}

	return pool, nil
}
