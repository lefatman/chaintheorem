// File: internal/persist/sessions_repo.go
package persist

import (
	"context"
	"database/sql"
)

type Session struct {
	Token     string
	UserID    int64
	ExpiresAt int64
	CreatedAt int64
}

type SessionsRepo struct {
	db      *sql.DB
	dialect Dialect
}

func NewSessionsRepo(db *sql.DB, dialect Dialect) *SessionsRepo {
	return &SessionsRepo{
		db:      db,
		dialect: dialect,
	}
}

func (r *SessionsRepo) Create(ctx context.Context, session Session) error {
	if r.db == nil {
		return ErrNilDB
	}
	_, err := r.db.ExecContext(ctx, insertSession(r.dialect), session.Token, session.UserID, session.ExpiresAt, session.CreatedAt)
	return err
}

func (r *SessionsRepo) Get(ctx context.Context, token string) (Session, error) {
	if r.db == nil {
		return Session{}, ErrNilDB
	}
	row := r.db.QueryRowContext(ctx, selectSessionByToken(r.dialect), token)
	var session Session
	if err := row.Scan(&session.Token, &session.UserID, &session.ExpiresAt, &session.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return Session{}, ErrNotFound
		}
		return Session{}, err
	}
	return session, nil
}

func (r *SessionsRepo) Delete(ctx context.Context, token string) error {
	if r.db == nil {
		return ErrNilDB
	}
	res, err := r.db.ExecContext(ctx, deleteSessionByToken(r.dialect), token)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func insertSession(dialect Dialect) string {
	if dialect == DialectPostgres {
		return `INSERT INTO sessions (token, user_id, expires_at, created_at) VALUES ($1, $2, $3, $4)`
	}
	return `INSERT INTO sessions (token, user_id, expires_at, created_at) VALUES (?, ?, ?, ?)`
}

func selectSessionByToken(dialect Dialect) string {
	if dialect == DialectPostgres {
		return `SELECT token, user_id, expires_at, created_at FROM sessions WHERE token = $1`
	}
	return `SELECT token, user_id, expires_at, created_at FROM sessions WHERE token = ?`
}

func deleteSessionByToken(dialect Dialect) string {
	if dialect == DialectPostgres {
		return `DELETE FROM sessions WHERE token = $1`
	}
	return `DELETE FROM sessions WHERE token = ?`
}
