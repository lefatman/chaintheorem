# STATE HANDOFF — Batch 03 (Persist)

## What this batch created / updated (scope-locked)
### Persistence layer
- `internal/persist/db.go`
- `internal/persist/migrate.go`
- `internal/persist/accounts_repo.go`
- `internal/persist/sessions_repo.go`
- `internal/persist/loadouts_repo.go`
- `internal/persist/progression_repo.go`
- `internal/persist/unlocks_repo.go`
  - Added stdlib-first DB open/ping with dialect validation and sentinel errors.
  - Added migration runner with version tracking and dialect-specific SQL.
  - Added CRUD repositories for accounts, sessions, loadouts, progression, and unlocks.

### Migrations
- `internal/persist/migrations/sqlite/001_init.sql`
- `internal/persist/migrations/postgres/001_init.sql`
  - Schema matches canon MVP tables with per-piece ability columns.

### Documentation updates
- `docs/ARCH_MAP/internal_persist.md`
- `docs/ARCH_MAP/README.md`
- `docs/DECISION_LEDGER.md`
- `docs/STATE_HANDOFF.md`

## Decisions appended
- DECISION 0007: Dialect-specific migration sets for SQLite/Postgres.
- DECISION 0008: Per-piece ability column naming for loadouts.
- DECISION 0009: Caller-owned driver registration.

## How to validate
1. `gofmt -w .`
2. `go test ./...`

## Next module to implement
- `docs/ARCH_MAP/internal_auth.md`
