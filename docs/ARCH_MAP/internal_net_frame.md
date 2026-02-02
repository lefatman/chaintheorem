# internal/net/frame

**Purpose:** Binary frame codec (u16 msg_type + u32 payload_len LE) used by ws_gateway.

## Canon inputs (authoritative)
- 03_PROTOCOL_CONTRACT — CONSOLIDATED
- PROJECT KERNEL v0.1 — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Strict bounds checks on `payload_len` (protect against OOM).
- Little-endian encoding/decoding.
- Unknown `msg_type` policy: drop + optionally send ERROR (per ws_gateway policy).

## Interfaces and boundaries
## Interfaces
- Used by `internal/ws_gateway` for all WSS reads/writes.

## File-by-file walkthrough (expected / required)
## Expected files
- `internal/net/frame/frame.go` — Encode/Decode header + payload, validation helpers
- (optional) `internal/net/frame/limits.go` — max frame size constants

## Gotchas / failure modes
## Hot-loop constraints
- No per-frame allocations beyond payload buffer (reuse where possible).
- Avoid reflection; avoid map iteration affecting output.

## Acceptance criteria
## Done when
- Fuzz/simple tests pass for malformed headers and length mismatches.
- Can round-trip encode/decode known frames.

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
- Only modify/create files under: `internal/net/frame/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
