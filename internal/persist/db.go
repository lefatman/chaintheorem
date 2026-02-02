// File: internal/persist/db.go
package persist

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Dialect string

const (
	DialectSQLite   Dialect = "sqlite"
	DialectPostgres Dialect = "postgres"
)

var (
	ErrNotFound = errors.New("persist: not found")
	ErrNilDB    = errors.New("persist: nil db")
)

type Config struct {
	Driver          string
	DSN             string
	Dialect         Dialect
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func (cfg Config) Validate() error {
	if cfg.Driver == "" {
		return fmt.Errorf("persist: driver required")
	}
	if cfg.DSN == "" {
		return fmt.Errorf("persist: dsn required")
	}
	switch cfg.Dialect {
	case DialectSQLite, DialectPostgres:
	default:
		return fmt.Errorf("persist: invalid dialect %q", cfg.Dialect)
	}
	return nil
}

func Open(ctx context.Context, cfg Config) (*sql.DB, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, err
	}
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func Ping(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return ErrNilDB
	}
	return db.PingContext(ctx)
}
