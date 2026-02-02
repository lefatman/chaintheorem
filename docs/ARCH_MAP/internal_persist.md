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

## File-by-file walkthrough (expected / required)
## Expected files
- `db.go` — open/health/ping
- `migrate.go` — migration runner
- `migrations/*.sql` — ordered schema files
- `accounts_repo.go`, `sessions_repo.go`, `loadouts_repo.go`, `progression_repo.go`, `unlocks_repo.go`

## Gotchas / failure modes
## Gotchas
- SQLite uses `INTEGER PRIMARY KEY` rowid; Postgres needs sequences/identity.
- Keep timestamps UTC integer seconds.

## Acceptance criteria
## Done when
- Fresh DB initializes via migrations on both SQLite and Postgres (or clearly separated dialect migrations).

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
- Only modify/create files under: `internal/persist/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
