# internal/ws_gateway

**Purpose:** WSS lifecycle, read/write loops, backpressure policy, session binding handshake.

## Canon inputs (authoritative)
- 03_PROTOCOL_CONTRACT — CONSOLIDATED
- PROJECT KERNEL v0.1 — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Two-loop design: read loop decodes frames; write loop drains bounded queue.
- Backpressure policy:
  - WORLD_DELTA can be dropped/coalesced.
  - Battle timelines/end must never be dropped (block or disconnect).
- Session binding via HELLO{token} before accepting game messages.

## Interfaces and boundaries
## Interfaces
- Upstream: net/http (HTTPS server)
- Downstream: `internal/net/frame`, `internal/router`
- Auth binding: calls `auth`/`persist` to validate token and attach player_id

## File-by-file walkthrough (expected / required)
## Expected files
- `server.go` — HTTP server integration + WS upgrade + accept loop
- `conn.go` — per-connection read/write loops + lifecycle
- `backpressure.go` — queue sizing, drop policies, classification
- (optional) `metrics.go` — counters/histograms hooks

## Gotchas / failure modes
## Gotchas
- Do not let write queue grow unbounded.
- Keep battle traffic reliable; prefer disconnect to lying.

## Acceptance criteria
## Done when
- Multiple clients can connect, HELLO, WELCOME, and exchange frames.
- Backpressure behaves as specified under artificial slow-client tests.

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
- Only modify/create files under: `internal/ws_gateway/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
