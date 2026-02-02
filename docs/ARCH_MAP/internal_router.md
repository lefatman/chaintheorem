# internal/router

**Purpose:** Dispatch frames to module handlers; schema-level validation only.

## Canon inputs (authoritative)
- 03_PROTOCOL_CONTRACT — CONSOLIDATED
- PROJECT KERNEL v0.1 — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Dispatch by msg_type.
- Validate schema invariants only (e.g., dx/dy in -1..1), not game legality.
- No heavy allocations; avoid map iteration for deterministic operations (use switch/array tables).

## Interfaces and boundaries
## Interfaces
- Upstream: ws_gateway passes (player_id, msg_type, payload)
- Downstream: module handlers: auth/world/aoi/chat/battle_mgr

## File-by-file walkthrough (expected / required)
## Expected files
- `router.go` — dispatch core
- `handlers.go` — module handler interfaces
- (optional) `errors.go` — ERROR message helpers

## Gotchas / failure modes
## Gotchas
- Router must not import heavy subsystems that cause cycles.
- Keep all msg_type handling in one obvious place (auditability).

## Acceptance criteria
## Done when
- Every msg_type in protocol has exactly one handler path or an explicit rejection.

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
- Only modify/create files under: `internal/router/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
