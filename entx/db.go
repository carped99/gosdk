package entx

import (
	"database/sql"
	"strings"
	"time"
)

func NewDB(cfg DatabaseConfig) (*sql.DB, error) {
	driverName := resolveDriverName(cfg.Driver)

	db, err := sql.Open(driverName, cfg.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	return db, nil
}

func resolveDriverName(driver string) string {
	switch strings.ToLower(driver) {
	case "postgres":
		return "pgx"
	default:
		return driver
	}
}
