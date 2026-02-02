// File: internal/auth/service.go
package auth

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"errors"
	"time"

	"example.com/mvp-repo/internal/persist"
)

var (
	ErrAccountExists      = errors.New("auth: account exists")
	ErrInvalidCredentials = errors.New("auth: invalid credentials")
	ErrInvalidToken       = errors.New("auth: invalid token")
	ErrTokenExpired       = errors.New("auth: token expired")
	ErrMissingIdentifiers = errors.New("auth: missing identifiers")
	ErrTokenTTLRequired   = errors.New("auth: token ttl required")
)

type Service struct {
	accounts *persist.AccountsRepo
	sessions *persist.SessionsRepo
	tokens   TokenGenerator
	tokenTTL time.Duration
	now      func() time.Time
}

type Config struct {
	TokenBytes int
	TokenTTL   time.Duration
	Now        func() time.Time
}

func NewService(accounts *persist.AccountsRepo, sessions *persist.SessionsRepo, cfg Config) (*Service, error) {
	if accounts == nil || sessions == nil {
		return nil, errors.New("auth: repos required")
	}
	if cfg.TokenTTL <= 0 {
		return nil, ErrTokenTTLRequired
	}
	tokens, err := NewTokenGenerator(cfg.TokenBytes)
	if err != nil {
		return nil, err
	}
	now := cfg.Now
	if now == nil {
		now = time.Now
	}
	return &Service{
		accounts: accounts,
		sessions: sessions,
		tokens:   tokens,
		tokenTTL: cfg.TokenTTL,
		now:      now,
	}, nil
}

func (s *Service) Register(ctx context.Context, email string, username string, password string) (persist.Session, error) {
	if email == "" || username == "" {
		return persist.Session{}, ErrMissingIdentifiers
	}
	if _, err := s.accounts.GetByEmail(ctx, email); err == nil {
		return persist.Session{}, ErrAccountExists
	} else if !errors.Is(err, persist.ErrNotFound) {
		return persist.Session{}, err
	}
	if _, err := s.accounts.GetByUsername(ctx, username); err == nil {
		return persist.Session{}, ErrAccountExists
	} else if !errors.Is(err, persist.ErrNotFound) {
		return persist.Session{}, err
	}
	hash, err := HashPassword(password)
	if err != nil {
		return persist.Session{}, err
	}
	now := s.now().Unix()
	account := persist.Account{
		Email:       email,
		Username:    username,
		PassHash:    []byte(hash),
		CreatedAt:   now,
		LastLoginAt: sql.NullInt64{},
	}
	userID, err := s.accounts.Create(ctx, account)
	if err != nil {
		return persist.Session{}, err
	}
	return s.newSession(ctx, userID)
}

func (s *Service) LoginByEmail(ctx context.Context, email string, password string) (persist.Session, error) {
	if email == "" {
		return persist.Session{}, ErrMissingIdentifiers
	}
	account, err := s.accounts.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, persist.ErrNotFound) {
			return persist.Session{}, ErrInvalidCredentials
		}
		return persist.Session{}, err
	}
	ok, err := VerifyPassword(password, string(account.PassHash))
	if err != nil {
		return persist.Session{}, err
	}
	if !ok {
		return persist.Session{}, ErrInvalidCredentials
	}
	now := s.now().Unix()
	if err := s.accounts.UpdateLastLogin(ctx, account.UserID, now); err != nil {
		return persist.Session{}, err
	}
	return s.newSession(ctx, account.UserID)
}

func (s *Service) LoginByUsername(ctx context.Context, username string, password string) (persist.Session, error) {
	if username == "" {
		return persist.Session{}, ErrMissingIdentifiers
	}
	account, err := s.accounts.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, persist.ErrNotFound) {
			return persist.Session{}, ErrInvalidCredentials
		}
		return persist.Session{}, err
	}
	ok, err := VerifyPassword(password, string(account.PassHash))
	if err != nil {
		return persist.Session{}, err
	}
	if !ok {
		return persist.Session{}, ErrInvalidCredentials
	}
	now := s.now().Unix()
	if err := s.accounts.UpdateLastLogin(ctx, account.UserID, now); err != nil {
		return persist.Session{}, err
	}
	return s.newSession(ctx, account.UserID)
}

func (s *Service) ValidateToken(ctx context.Context, token string) (persist.Session, error) {
	if token == "" {
		return persist.Session{}, ErrInvalidToken
	}
	session, err := s.sessions.Get(ctx, token)
	if err != nil {
		if errors.Is(err, persist.ErrNotFound) {
			return persist.Session{}, ErrInvalidToken
		}
		return persist.Session{}, err
	}
	if subtle.ConstantTimeCompare([]byte(session.Token), []byte(token)) != 1 {
		return persist.Session{}, ErrInvalidToken
	}
	now := s.now().Unix()
	if session.ExpiresAt <= now {
		_ = s.sessions.Delete(ctx, token)
		return persist.Session{}, ErrTokenExpired
	}
	return session, nil
}

func (s *Service) RevokeToken(ctx context.Context, token string) error {
	if token == "" {
		return ErrInvalidToken
	}
	if err := s.sessions.Delete(ctx, token); err != nil {
		if errors.Is(err, persist.ErrNotFound) {
			return ErrInvalidToken
		}
		return err
	}
	return nil
}

func (s *Service) newSession(ctx context.Context, userID int64) (persist.Session, error) {
	token, err := s.tokens.NewToken()
	if err != nil {
		return persist.Session{}, err
	}
	now := s.now().Unix()
	session := persist.Session{
		Token:     token,
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now + int64(s.tokenTTL.Seconds()),
	}
	if err := s.sessions.Create(ctx, session); err != nil {
		return persist.Session{}, err
	}
	return session, nil
}
