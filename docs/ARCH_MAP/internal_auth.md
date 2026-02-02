---
owner: internal/auth
status: DONE
generated_files:
  - internal/auth/service.go
  - internal/auth/password.go
  - internal/auth/tokens.go
touchpoints:
  - docs/DECISION_LEDGER.md
  - docs/ARCH_MAP/README.md
  - docs/STATE_HANDOFF.md
  - go.mod
  - go.sum
last_updated: 2026-02-02
---

# internal/auth

**Purpose:** Register/login/reset + session tokens, password hashing.

## Canon inputs (authoritative)
- 06_ACCOUNTS_PERSISTENCE_SOCIAL — CONSOLIDATED
- PROJECT KERNEL v0.1 — CONSOLIDATED
- MVP System Architecture — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- HTTPS register/login/reset flows; WSS uses HELLO token binding.
- Password hashing via bcrypt or argon2id (decision is ledger-owned).
- Session tokens are opaque random bytes stored server-side with expiry (MVP).

## Interfaces and boundaries
## Interfaces
- Uses `internal/persist` for accounts/sessions.
- Exposed via `internal/httpapi` handlers.

## File-by-file walkthrough (expected / required)
## Expected files
- `service.go` — high-level auth operations
- `password.go` — hash + verify
- `tokens.go` — token generation + expiry
- (optional) `ratelimit.go` — basic anti-abuse

## Gotchas / failure modes
## Gotchas
- Never log raw passwords or session tokens.
- Token comparison should be constant-time if feasible.

## Acceptance criteria
## Done when
- Register/login works end-to-end; HELLO binds a WSS connection using issued token.

## Generated/Modified Files
- `internal/auth/service.go`
- `internal/auth/password.go`
- `internal/auth/tokens.go`

## Interfaces / Contracts
- `Service` with register/login/token validation for `internal/httpapi`.
- Uses `internal/persist.AccountsRepo` and `internal/persist.SessionsRepo`.
- Password hashing uses argon2id encoded hash strings.

## Algorithmic Invariants Implemented
- Opaque session tokens are random bytes encoded for transport and stored with expiry.
- Token validation performs constant-time comparisons when possible.
- Password verification uses constant-time hash comparison.

## Remaining Work
- None.

### Prompt seed for this subdirectory (for later)
Use this as the nucleus for a per-subdir generator prompt.

**Required attachments**
- `docs/CANON_LOCK.md`
- `docs/DECISION_LEDGER.md`
- `docs/STATE_HANDOFF.md` (latest)
- `TREE.txt`
- Any existing files under this subdir
- The governing design docs for this subdir (see “Canon inputs” above)

**Scope lock**
- Only modify/create files under: `internal/auth/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
