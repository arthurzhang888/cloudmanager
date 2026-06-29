package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresDB wraps a PostgreSQL connection pool
type PostgresDB struct {
	Pool *pgxpool.Pool
}

// NewPostgresDB creates a new PostgreSQL connection pool
func NewPostgresDB(connString string) (*PostgresDB, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{Pool: pool}, nil
}

// Close closes the database connection pool
func (db *PostgresDB) Close() {
	db.Pool.Close()
}

// Health checks if the database is reachable
func (db *PostgresDB) Health(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

// Exec executes a query without returning rows
func (db *PostgresDB) Exec(ctx context.Context, sql string, args ...interface{}) (interface{}, error) {
	return db.Pool.Exec(ctx, sql, args...)
}

// Query executes a query that returns multiple rows
func (db *PostgresDB) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	return db.Pool.Query(ctx, sql, args...)
}

// QueryRow executes a query that returns a single row
func (db *PostgresDB) QueryRow(ctx context.Context, sql string, args ...interface{}) Row {
	return db.Pool.QueryRow(ctx, sql, args...)
}

// Rows interface for database rows
type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close()
	Err() error
}

// Row interface for a single database row
type Row interface {
	Scan(dest ...interface{}) error
}
