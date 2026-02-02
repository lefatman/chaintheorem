// File: internal/auth/password.go
package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	passwordSaltBytes = 16
	passwordHashBytes = 32
	passwordTime      = 3
	passwordMemory    = 64 * 1024
	passwordThreads   = 2
)

var (
	ErrInvalidPassword = errors.New("auth: invalid password")
	ErrInvalidHash     = errors.New("auth: invalid password hash")
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrInvalidPassword
	}
	salt := make([]byte, passwordSaltBytes)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, passwordTime, passwordMemory, passwordThreads, passwordHashBytes)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", passwordMemory, passwordTime, passwordThreads, b64Salt, b64Hash)
	return encoded, nil
}

func VerifyPassword(password string, encoded string) (bool, error) {
	if password == "" || encoded == "" {
		return false, ErrInvalidPassword
	}
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != "v=19" {
		return false, ErrInvalidHash
	}
	mem, timeCost, threads, err := parseArgon2Params(parts[3])
	if err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, ErrInvalidHash
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, ErrInvalidHash
	}
	sum := argon2.IDKey([]byte(password), salt, timeCost, mem, threads, uint32(len(hash)))
	if subtle.ConstantTimeCompare(sum, hash) != 1 {
		return false, nil
	}
	return true, nil
}

func parseArgon2Params(raw string) (uint32, uint32, uint8, error) {
	parts := strings.Split(raw, ",")
	if len(parts) != 3 {
		return 0, 0, 0, ErrInvalidHash
	}
	mem, err := parseArgon2Param(parts[0], "m=")
	if err != nil {
		return 0, 0, 0, err
	}
	timeCost, err := parseArgon2Param(parts[1], "t=")
	if err != nil {
		return 0, 0, 0, err
	}
	threadsVal, err := parseArgon2Param(parts[2], "p=")
	if err != nil {
		return 0, 0, 0, err
	}
	if threadsVal > 255 {
		return 0, 0, 0, ErrInvalidHash
	}
	return uint32(mem), uint32(timeCost), uint8(threadsVal), nil
}

func parseArgon2Param(part string, prefix string) (int, error) {
	if !strings.HasPrefix(part, prefix) {
		return 0, ErrInvalidHash
	}
	val, err := strconv.Atoi(strings.TrimPrefix(part, prefix))
	if err != nil || val <= 0 {
		return 0, ErrInvalidHash
	}
	return val, nil
}
