# STATE HANDOFF — Batch 04 (Auth)

## What this batch created / updated (scope-locked)
### Auth service
- `internal/auth/service.go`
- `internal/auth/password.go`
- `internal/auth/tokens.go`
  - Added auth service for register/login/token validation using persist repos.
  - Added argon2id password hashing and constant-time verification.
  - Added opaque random token generation with expiry.

### Documentation updates
- `docs/ARCH_MAP/internal_auth.md`
- `docs/ARCH_MAP/README.md`
- `docs/DECISION_LEDGER.md`
- `docs/STATE_HANDOFF.md`

## Decisions appended
- DECISION 0010: Password hashing uses argon2id.

## How to validate
1. `gofmt -w .`
2. `go test ./...`

## Next module to implement
- `docs/ARCH_MAP/internal_httpapi.md`
