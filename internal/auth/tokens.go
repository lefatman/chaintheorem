// File: internal/auth/tokens.go
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

const (
	DefaultTokenBytes = 32
	MinTokenBytes     = 16
)

var ErrInvalidTokenSize = errors.New("auth: invalid token size")

type TokenGenerator struct {
	size int
}

func NewTokenGenerator(size int) (TokenGenerator, error) {
	if size == 0 {
		size = DefaultTokenBytes
	}
	if size < MinTokenBytes {
		return TokenGenerator{}, ErrInvalidTokenSize
	}
	return TokenGenerator{size: size}, nil
}

func (g TokenGenerator) NewToken() (string, error) {
	buf := make([]byte, g.size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
