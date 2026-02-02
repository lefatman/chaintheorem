// File: internal/persist/migrate.go
package persist

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/sqlite/*.sql
var sqliteMigrations embed.FS

//go:embed migrations/postgres/*.sql
var postgresMigrations embed.FS

type migration struct {
	version int
	name    string
	sql     string
}

func Migrate(ctx context.Context, db *sql.DB, dialect Dialect) error {
	if db == nil {
		return ErrNilDB
	}
	migrations, err := loadMigrations(dialect)
	if err != nil {
		return err
	}
	if err := ensureMigrationsTable(ctx, db); err != nil {
		return err
	}
	applied, err := loadAppliedMigrations(ctx, db)
	if err != nil {
		return err
	}
	insertSQL := insertMigrationSQL(dialect)
	for _, m := range migrations {
		if applied[m.version] {
			continue
		}
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, m.sql); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("persist: migrate %s: %w", m.name, err)
		}
		if _, err := tx.ExecContext(ctx, insertSQL, m.version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("persist: record migration %s: %w", m.name, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("persist: commit migration %s: %w", m.name, err)
		}
	}
	return nil
}

func ensureMigrationsTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER PRIMARY KEY)`)
	return err
}

func loadAppliedMigrations(ctx context.Context, db *sql.DB) (map[int]bool, error) {
	rows, err := db.QueryContext(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return applied, nil
}

func insertMigrationSQL(dialect Dialect) string {
	switch dialect {
	case DialectPostgres:
		return `INSERT INTO schema_migrations(version) VALUES ($1)`
	default:
		return `INSERT INTO schema_migrations(version) VALUES (?)`
	}
}

func loadMigrations(dialect Dialect) ([]migration, error) {
	var root fs.FS
	var glob string
	switch dialect {
	case DialectSQLite:
		root = sqliteMigrations
		glob = "migrations/sqlite/*.sql"
	case DialectPostgres:
		root = postgresMigrations
		glob = "migrations/postgres/*.sql"
	default:
		return nil, fmt.Errorf("persist: unsupported dialect %q", dialect)
	}
	paths, err := fs.Glob(root, glob)
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)
	migrations := make([]migration, 0, len(paths))
	for _, path := range paths {
		raw, err := fs.ReadFile(root, path)
		if err != nil {
			return nil, err
		}
		name := strings.TrimPrefix(path, "migrations/")
		version, err := parseMigrationVersion(path)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, migration{
			version: version,
			name:    name,
			sql:     string(raw),
		})
	}
	return migrations, nil
}

func parseMigrationVersion(path string) (int, error) {
	base := path[strings.LastIndex(path, "/")+1:]
	parts := strings.SplitN(base, "_", 2)
	if len(parts) < 1 {
		return 0, fmt.Errorf("persist: invalid migration name %q", path)
	}
	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("persist: invalid migration version in %q", path)
	}
	return version, nil
}
