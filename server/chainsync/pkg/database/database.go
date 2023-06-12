package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
)

func InitDatabase(dsn string) error {
	var err error
	pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Set connection pool configurations
	config, err := pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("failed to acquire connection from the pool: %w", err)
	}
	defer config.Release()

	if err != nil {
		return fmt.Errorf("failed to configure the database connection: %w", err)
	}

	return nil
}

func CloseDatabase() {
	if pool != nil {
		pool.Close()
	}
}

func GetConnection() (*pgxpool.Conn, error) {
	if pool == nil {
		return nil, fmt.Errorf("database connection pool is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection from the pool: %w", err)
	}

	return conn, nil
}
