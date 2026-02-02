# internal/app

**Purpose:** Dependency wiring + server bootstrap composition.

## Canon inputs (authoritative)
- MVP System Architecture — CONSOLIDATED
- PROJECT KERNEL v0.1 — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Composition root only: wire modules and start servers.
- No business logic beyond configuration and lifecycle.

## Interfaces and boundaries
## Interfaces
- Called from `cmd/server/main.go`.

## File-by-file walkthrough (expected / required)
## Expected files
- `app.go` — wire dependencies and expose Start/Stop
- (optional) `wiring.go` — small helpers to keep app.go short

## Gotchas / failure modes
## Gotchas
- Don’t hide decisions here; put them in configs or ledger.

## Acceptance criteria
## Done when
- `cmd/server` can start the server with this app wiring.

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
- Only modify/create files under: `internal/app/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
