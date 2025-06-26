package entx

import (
	"fmt"
	"time"
)

type DatabaseDriver string

const (
	DatabaseDriverPostgres DatabaseDriver = "postgres"
	DatabaseDriverMySQL    DatabaseDriver = "mysql"
	DatabaseDriverSQLite   DatabaseDriver = "sqlite"
	DatabaseDriverMSSQL    DatabaseDriver = "mssql"
)

type DatabaseConfig struct {
	Driver         DatabaseDriver `koanf:"driver"`
	DSN            string         `koanf:"dsn"`
	ConnectTimeout time.Duration  `koanf:"minIdleConns"`
	MinIdleConns   int            `koanf:"minIdleConns"`
	MaxIdleConns   int            `koanf:"maxIdleConns"`
	MaxOpenConns   int            `koanf:"maxOpenConns"`
	MaxLifetime    time.Duration  `koanf:"connMaxLifetime"`
	MaxIdleTime    time.Duration  `koanf:"connMaxLifetime"`
}

func (c DatabaseConfig) validate() error {
	if c.Driver == "" {
		return fmt.Errorf("database driver cannot be empty")
	}

	if c.DSN == "" {
		return fmt.Errorf("database DSN cannot be empty")
	}
	if c.MinIdleConns < 0 {
		return fmt.Errorf("database pool min idle connections cannot be negative")
	}

	if c.MaxOpenConns < 0 {
		return fmt.Errorf("database pool max open connections cannot be negative")
	}

	if c.MaxLifetime < 0 {
		return fmt.Errorf("database pool connection max lifetime cannot be negative")
	}

	if c.MinIdleConns > c.MaxOpenConns {
		return fmt.Errorf("database pool min idle connections cannot be greater than max open connections")
	}

	return nil
}

var DefaultDatabaseConfig = DatabaseConfig{
	Driver: "postgres",
	DSN:    "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
}
