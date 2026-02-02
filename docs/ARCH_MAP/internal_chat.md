# internal/chat

**Purpose:** Global chat broadcast + rate limiting.

## Canon inputs (authoritative)
- PROJECT KERNEL v0.1 — CONSOLIDATED
- 06_ACCOUNTS_PERSISTENCE_SOCIAL — CONSOLIDATED
- 03_PROTOCOL_CONTRACT — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- MVP: global channel only.
- Rate limit per user; optional persistence deferred.

## Interfaces and boundaries
## Interfaces
- Router receives CHAT_SEND and forwards to chat service.
- Chat service sends CHAT_EVENT to all connected sessions.

## File-by-file walkthrough (expected / required)
## Expected files
- `service.go` — broadcast logic
- `ratelimit.go` — token bucket

## Gotchas / failure modes
## Gotchas
- Never block the world tick on chat broadcast; use bounded queues.

## Acceptance criteria
## Done when
- Chat messages broadcast to all online clients with rate limiting.

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
- Only modify/create files under: `internal/chat/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
