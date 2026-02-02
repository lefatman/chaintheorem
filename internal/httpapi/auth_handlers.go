// File: internal/httpapi/auth_handlers.go
package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"example.com/mvp-repo/internal/auth"
	"example.com/mvp-repo/internal/persist"
)

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type resetRequest struct {
	Email string `json:"email"`
}

type sessionResponse struct {
	Token     string `json:"token"`
	UserID    int64  `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
	ExpiresAt int64  `json:"expires_at"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	var req registerRequest
	if err := s.decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid json payload")
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.TrimSpace(req.Username)
	if req.Email == "" || req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "missing_fields", "email, username, and password are required")
		return
	}
	session, err := s.auth.Register(r.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrAccountExists):
			writeError(w, http.StatusConflict, "account_exists", "account already exists")
		case errors.Is(err, auth.ErrMissingIdentifiers):
			writeError(w, http.StatusBadRequest, "missing_fields", "email and username are required")
		default:
			writeError(w, http.StatusInternalServerError, "server_error", "unable to register")
		}
		return
	}
	writeData(w, http.StatusCreated, sessionResponse{
		Token:     session.Token,
		UserID:    session.UserID,
		CreatedAt: session.CreatedAt,
		ExpiresAt: session.ExpiresAt,
	})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	var req loginRequest
	if err := s.decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid json payload")
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.TrimSpace(req.Username)
	if req.Password == "" {
		writeError(w, http.StatusBadRequest, "missing_fields", "password is required")
		return
	}
	var (
		session persist.Session
		err     error
	)
	switch {
	case req.Email != "" && req.Username != "":
		writeError(w, http.StatusBadRequest, "invalid_credentials", "provide only email or username")
		return
	case req.Email != "":
		session, err = s.auth.LoginByEmail(r.Context(), req.Email, req.Password)
	case req.Username != "":
		session, err = s.auth.LoginByUsername(r.Context(), req.Username, req.Password)
	default:
		writeError(w, http.StatusBadRequest, "missing_fields", "email or username is required")
		return
	}
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			writeError(w, http.StatusUnauthorized, "invalid_credentials", "invalid credentials")
		case errors.Is(err, auth.ErrMissingIdentifiers):
			writeError(w, http.StatusBadRequest, "missing_fields", "email or username is required")
		default:
			writeError(w, http.StatusInternalServerError, "server_error", "unable to login")
		}
		return
	}
	writeData(w, http.StatusOK, sessionResponse{
		Token:     session.Token,
		UserID:    session.UserID,
		CreatedAt: session.CreatedAt,
		ExpiresAt: session.ExpiresAt,
	})
}

func (s *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	var req resetRequest
	if err := s.decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid json payload")
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "missing_fields", "email is required")
		return
	}
	if err := s.auth.ResetPassword(r.Context(), req.Email); err != nil {
		writeError(w, http.StatusInternalServerError, "server_error", "unable to reset password")
		return
	}
	writeData(w, http.StatusOK, map[string]string{
		"status": "reset_requested",
	})
}
