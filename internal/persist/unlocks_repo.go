// File: internal/persist/unlocks_repo.go
package persist

import (
	"context"
	"database/sql"
)

type UserUnlock struct {
	UserID     int64
	FlagID     int64
	UnlockedAt int64
}

type UnlocksRepo struct {
	db      *sql.DB
	dialect Dialect
}

func NewUnlocksRepo(db *sql.DB, dialect Dialect) *UnlocksRepo {
	return &UnlocksRepo{
		db:      db,
		dialect: dialect,
	}
}

func (r *UnlocksRepo) Add(ctx context.Context, unlock UserUnlock) error {
	if r.db == nil {
		return ErrNilDB
	}
	_, err := r.db.ExecContext(ctx, insertUnlock(r.dialect), unlock.UserID, unlock.FlagID, unlock.UnlockedAt)
	return err
}

func (r *UnlocksRepo) ListByUser(ctx context.Context, userID int64) ([]UserUnlock, error) {
	if r.db == nil {
		return nil, ErrNilDB
	}
	rows, err := r.db.QueryContext(ctx, selectUnlocksByUser(r.dialect), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	unlocks := make([]UserUnlock, 0)
	for rows.Next() {
		var unlock UserUnlock
		if err := rows.Scan(&unlock.UserID, &unlock.FlagID, &unlock.UnlockedAt); err != nil {
			return nil, err
		}
		unlocks = append(unlocks, unlock)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return unlocks, nil
}

func insertUnlock(dialect Dialect) string {
	if dialect == DialectPostgres {
		return `INSERT INTO user_unlocks (user_id, flag_id, unlocked_at) VALUES ($1, $2, $3) ON CONFLICT (user_id, flag_id) DO NOTHING`
	}
	return `INSERT INTO user_unlocks (user_id, flag_id, unlocked_at) VALUES (?, ?, ?) ON CONFLICT(user_id, flag_id) DO NOTHING`
}

func selectUnlocksByUser(dialect Dialect) string {
	if dialect == DialectPostgres {
		return `SELECT user_id, flag_id, unlocked_at FROM user_unlocks WHERE user_id = $1 ORDER BY flag_id`
	}
	return `SELECT user_id, flag_id, unlocked_at FROM user_unlocks WHERE user_id = ? ORDER BY flag_id`
}
