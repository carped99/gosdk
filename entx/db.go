package entx

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// DB represents a database connection with extended functionality
type DB struct {
	*sql.DB
	config DatabaseConfig
}

// NewDB creates a new database connection with the given configuration.
// It applies all the necessary settings and validates the connection.
func NewDB(cfg DatabaseConfig) (*DB, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err)
	}

	driverName := resolveDriverName(cfg.Driver)
	db, err := sql.Open(driverName, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Apply pool settings
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)
	db.SetConnMaxIdleTime(cfg.MaxIdleTime)

	// Validate connection
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		DB:     db,
		config: cfg,
	}, nil
}

// resolveDriverName returns the appropriate driver name for the given database type.
// It handles special cases and driver-specific requirements.
func resolveDriverName(driver DatabaseDriver) string {
	switch strings.ToLower(string(driver)) {
	case string(DatabaseDriverPostgres):
		return "pgx"
	case string(DatabaseDriverMySQL):
		return "mysql"
	default:
		return string(driver)
	}
}

// Close closes the database connection and releases all resources.
func (db *DB) Close() error {
	return db.DB.Close()
}

// Config returns the database configuration.
func (db *DB) Config() DatabaseConfig {
	return db.config
}

// BeginTx starts a new transaction with the given context and options.
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.DB.BeginTx(ctx, opts)
}

// ExecContext executes a query without returning any rows.
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.DB.ExecContext(ctx, query, args...)
}

// QueryContext executes a query that returns rows.
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.DB.QueryContext(ctx, query, args...)
}

// QueryRowContext executes a query that returns a single row.
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.DB.QueryRowContext(ctx, query, args...)
}
