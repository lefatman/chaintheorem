---
status: done
owner: internal/persist
generated_files:
  - internal/persist/db.go
  - internal/persist/migrate.go
  - internal/persist/accounts_repo.go
  - internal/persist/sessions_repo.go
  - internal/persist/loadouts_repo.go
  - internal/persist/progression_repo.go
  - internal/persist/unlocks_repo.go
  - internal/persist/migrations/sqlite/001_init.sql
  - internal/persist/migrations/postgres/001_init.sql
touchpoints:
  - docs/DECISION_LEDGER.md
  - docs/ARCH_MAP/README.md
  - docs/STATE_HANDOFF.md
depends_on: []
last_updated: 2026-02-02
---

# internal/persist

**Purpose:** DB schema, migrations, repositories; SQLite dev + Postgres prod.

## Canon inputs (authoritative)
- 06_ACCOUNTS_PERSISTENCE_SOCIAL — CONSOLIDATED
- PROJECT KERNEL v0.1 — CONSOLIDATED
- MVP System Architecture — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- SQLite dev + Postgres prod; schema portability is required.
- Migration runner: version table + ordered migrations.
- Repos expose minimal CRUD; business rules stay outside persist.

## Interfaces and boundaries
## Canon schema (MVP extract)
```sql
accounts(user_id PK, email UNIQUE, username UNIQUE, pass_hash, created_at, last_login_at)
sessions(token PK, user_id, expires_at, created_at)
army_loadouts(user_id PK, element_id, army_ability_1..4, ability_*piece, item_1..4, updated_at)
progression(user_id PK, level, xp)
user_unlocks(user_id, flag_id, unlocked_at, PK(user_id, flag_id))
```

## Generated/Modified Files
- `internal/persist/db.go`
- `internal/persist/migrate.go`
- `internal/persist/accounts_repo.go`
- `internal/persist/sessions_repo.go`
- `internal/persist/loadouts_repo.go`
- `internal/persist/progression_repo.go`
- `internal/persist/unlocks_repo.go`
- `internal/persist/migrations/sqlite/001_init.sql`
- `internal/persist/migrations/postgres/001_init.sql`

## Interfaces / Contracts
- `persist.Config` + `persist.Open(ctx, cfg)` + `persist.Ping(ctx, db)`
- `persist.Migrate(ctx, db, dialect)` with per-dialect embedded migrations
- `AccountsRepo`, `SessionsRepo`, `LoadoutsRepo`, `ProgressionRepo`, `UnlocksRepo` CRUD helpers
- Sentinel errors: `ErrNotFound`, `ErrNilDB`

## Algorithmic Invariants Implemented
- Ordered, versioned migrations tracked in `schema_migrations`.
- Separate SQLite/Postgres migrations to preserve type correctness while keeping schema parity.
- Loadout storage uses fixed columns for army slots, per-piece assignments, and item slots.
- No business rules inside persistence layer (CRUD-only).

## Remaining Work
- None.
