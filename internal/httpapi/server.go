// File: internal/httpapi/server.go
package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"example.com/mvp-repo/internal/persist"
)

const defaultMaxBodyBytes int64 = 1 << 20

var (
	ErrAuthServiceRequired    = errors.New("httpapi: auth service required")
	ErrLoadoutServiceRequired = errors.New("httpapi: loadout service required")
	ErrListenAddrRequired     = errors.New("httpapi: listen addr required")
)

type Config struct {
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxBodyBytes      int64
}

type Server struct {
	auth         AuthService
	loadouts     LoadoutService
	maxBodyBytes int64
	mux          *http.ServeMux
	server       *http.Server
}

type AuthService interface {
	Register(ctx context.Context, email string, username string, password string) (persist.Session, error)
	LoginByEmail(ctx context.Context, email string, password string) (persist.Session, error)
	LoginByUsername(ctx context.Context, username string, password string) (persist.Session, error)
	ResetPassword(ctx context.Context, email string) error
	ValidateToken(ctx context.Context, token string) (persist.Session, error)
}

type LoadoutService interface {
	Get(ctx context.Context, userID int64) (persist.Loadout, error)
	Update(ctx context.Context, userID int64, input LoadoutInput) (persist.Loadout, error)
}

func NewServer(cfg Config, auth AuthService, loadouts LoadoutService) (*Server, error) {
	if auth == nil {
		return nil, ErrAuthServiceRequired
	}
	if loadouts == nil {
		return nil, ErrLoadoutServiceRequired
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 5 * time.Second
	}
	if cfg.ReadHeaderTimeout == 0 {
		cfg.ReadHeaderTimeout = 5 * time.Second
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 10 * time.Second
	}
	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = 30 * time.Second
	}
	if cfg.MaxBodyBytes <= 0 {
		cfg.MaxBodyBytes = defaultMaxBodyBytes
	}
	mux := http.NewServeMux()
	s := &Server{
		auth:         auth,
		loadouts:     loadouts,
		maxBodyBytes: cfg.MaxBodyBytes,
		mux:          mux,
	}
	mux.HandleFunc("/api/auth/register", s.handleRegister)
	mux.HandleFunc("/api/auth/login", s.handleLogin)
	mux.HandleFunc("/api/auth/reset", s.handleReset)
	mux.HandleFunc("/api/loadout", s.handleLoadout)
	server := &http.Server{
		Handler:           mux,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
	s.server = server
	return s, nil
}

func (s *Server) ListenAndServe(addr string) error {
	if addr == "" {
		return ErrListenAddrRequired
	}
	s.server.Addr = addr
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type apiResponse struct {
	OK    bool      `json:"ok"`
	Data  any       `json:"data,omitempty"`
	Error *apiError `json:"error,omitempty"`
}

func (s *Server) decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, s.maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("httpapi: unexpected trailing data")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload apiResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	_ = enc.Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, apiResponse{
		OK: false,
		Error: &apiError{
			Code:    code,
			Message: message,
		},
	})
}

func writeData(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, apiResponse{
		OK:   true,
		Data: data,
	})
}
