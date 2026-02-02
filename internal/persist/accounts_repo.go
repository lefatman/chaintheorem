// File: internal/persist/accounts_repo.go
package persist

import (
	"context"
	"database/sql"
)

type Account struct {
	UserID      int64
	Email       string
	Username    string
	PassHash    []byte
	CreatedAt   int64
	LastLoginAt sql.NullInt64
}

type AccountsRepo struct {
	db      *sql.DB
	dialect Dialect
}

func NewAccountsRepo(db *sql.DB, dialect Dialect) *AccountsRepo {
	return &AccountsRepo{
		db:      db,
		dialect: dialect,
	}
}

func (r *AccountsRepo) Create(ctx context.Context, account Account) (int64, error) {
	if r.db == nil {
		return 0, ErrNilDB
	}
	switch r.dialect {
	case DialectPostgres:
		row := r.db.QueryRowContext(ctx, createAccountPostgres, account.Email, account.Username, account.PassHash, account.CreatedAt, account.LastLoginAt)
		var userID int64
		if err := row.Scan(&userID); err != nil {
			return 0, err
		}
		return userID, nil
	default:
		res, err := r.db.ExecContext(ctx, createAccountSQLite, account.Email, account.Username, account.PassHash, account.CreatedAt, account.LastLoginAt)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}
}

func (r *AccountsRepo) GetByID(ctx context.Context, userID int64) (Account, error) {
	return r.getSingle(ctx, selectAccountByID(r.dialect), userID)
}

func (r *AccountsRepo) GetByEmail(ctx context.Context, email string) (Account, error) {
	return r.getSingle(ctx, selectAccountByEmail(r.dialect), email)
}

func (r *AccountsRepo) GetByUsername(ctx context.Context, username string) (Account, error) {
	return r.getSingle(ctx, selectAccountByUsername(r.dialect), username)
}

func (r *AccountsRepo) UpdateLastLogin(ctx context.Context, userID int64, lastLoginAt int64) error {
	if r.db == nil {
		return ErrNilDB
	}
	query := updateAccountLastLogin(r.dialect)
	res, err := r.db.ExecContext(ctx, query, lastLoginAt, userID)
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

func (r *AccountsRepo) getSingle(ctx context.Context, query string, arg any) (Account, error) {
	if r.db == nil {
		return Account{}, ErrNilDB
	}
	row := r.db.QueryRowContext(ctx, query, arg)
	var account Account
	if err := row.Scan(&account.UserID, &account.Email, &account.Username, &account.PassHash, &account.CreatedAt, &account.LastLoginAt); err != nil {
		if err == sql.ErrNoRows {
			return Account{}, ErrNotFound
		}
		return Account{}, err
	}
	return account, nil
}

func selectAccountByID(dialect Dialect) string {
	if dialect == DialectPostgres {
		return `SELECT user_id, email, username, pass_hash, created_at, last_login_at FROM accounts WHERE user_id = $1`
	}
	return `SELECT user_id, email, username, pass_hash, created_at, last_login_at FROM accounts WHERE user_id = ?`
}

func selectAccountByEmail(dialect Dialect) string {
	if dialect == DialectPostgres {
		return `SELECT user_id, email, username, pass_hash, created_at, last_login_at FROM accounts WHERE email = $1`
	}
	return `SELECT user_id, email, username, pass_hash, created_at, last_login_at FROM accounts WHERE email = ?`
}

func selectAccountByUsername(dialect Dialect) string {
	if dialect == DialectPostgres {
		return `SELECT user_id, email, username, pass_hash, created_at, last_login_at FROM accounts WHERE username = $1`
	}
	return `SELECT user_id, email, username, pass_hash, created_at, last_login_at FROM accounts WHERE username = ?`
}

func updateAccountLastLogin(dialect Dialect) string {
	if dialect == DialectPostgres {
		return `UPDATE accounts SET last_login_at = $1 WHERE user_id = $2`
	}
	return `UPDATE accounts SET last_login_at = ? WHERE user_id = ?`
}

const createAccountSQLite = `INSERT INTO accounts (email, username, pass_hash, created_at, last_login_at) VALUES (?, ?, ?, ?, ?)`
const createAccountPostgres = `INSERT INTO accounts (email, username, pass_hash, created_at, last_login_at) VALUES ($1, $2, $3, $4, $5) RETURNING user_id`
