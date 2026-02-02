// File: internal/persist/progression_repo.go
package persist

import (
	"context"
	"database/sql"
)

type Progression struct {
	UserID int64
	Level  int64
	XP     int64
}

type ProgressionRepo struct {
	db      *sql.DB
	dialect Dialect
}

func NewProgressionRepo(db *sql.DB, dialect Dialect) *ProgressionRepo {
	return &ProgressionRepo{
		db:      db,
		dialect: dialect,
	}
}

func (r *ProgressionRepo) Get(ctx context.Context, userID int64) (Progression, error) {
	if r.db == nil {
		return Progression{}, ErrNilDB
	}
	row := r.db.QueryRowContext(ctx, selectProgression(r.dialect), userID)
	var progression Progression
	if err := row.Scan(&progression.UserID, &progression.Level, &progression.XP); err != nil {
		if err == sql.ErrNoRows {
			return Progression{}, ErrNotFound
		}
		return Progression{}, err
	}
	return progression, nil
}

func (r *ProgressionRepo) Upsert(ctx context.Context, progression Progression) error {
	if r.db == nil {
		return ErrNilDB
	}
	_, err := r.db.ExecContext(ctx, upsertProgression(r.dialect), progression.UserID, progression.Level, progression.XP)
	return err
}

func selectProgression(dialect Dialect) string {
	if dialect == DialectPostgres {
		return `SELECT user_id, level, xp FROM progression WHERE user_id = $1`
	}
	return `SELECT user_id, level, xp FROM progression WHERE user_id = ?`
}

func upsertProgression(dialect Dialect) string {
	if dialect == DialectPostgres {
		return `INSERT INTO progression (user_id, level, xp) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET level = EXCLUDED.level, xp = EXCLUDED.xp`
	}
	return `INSERT INTO progression (user_id, level, xp) VALUES (?, ?, ?) ON CONFLICT(user_id) DO UPDATE SET level = excluded.level, xp = excluded.xp`
}
