# internal/httpapi

**Purpose:** HTTPS JSON endpoints for auth + loadout editing + dev hooks.

## Canon inputs (authoritative)
- 06_ACCOUNTS_PERSISTENCE_SOCIAL — CONSOLIDATED
- PROJECT KERNEL v0.1 — CONSOLIDATED
- MVP System Architecture — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- JSON over HTTPS is allowed (not hot path).
- Keep handlers thin: validate, call service, return response.

## Interfaces and boundaries
## Interfaces
- Uses `auth`, `persist`, `loadout`, `battle_mgr` (for dev hooks).

## File-by-file walkthrough (expected / required)
## Expected files
- `server.go` — mux/routes
- `auth_handlers.go` — register/login/reset
- `loadout_handlers.go` — GET/POST loadout (for MVP testing)
- (optional) `dev_handlers.go` — dev-only battle start hooks

## Gotchas / failure modes
## Gotchas
- Keep CORS and cookie/token policies consistent; document choices in ledger.

## Acceptance criteria
## Done when
- You can create an account and obtain a session token via HTTPS endpoints.

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
- Only modify/create files under: `internal/httpapi/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
