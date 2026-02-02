# internal/aoi

**Purpose:** Grid AOI, watcher sets, diff assembly (snapshot/delta), resync.

## Canon inputs (authoritative)
- 04_AOI_AND_REPLICATION — CONSOLIDATED
- MVP System Architecture — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Grid AOI: entities bucketed by cell.
- Per-player watcher sets: square neighborhood radius R cells.
- Diff production scans only watcher cells; stable sort by entity_id.
- Backpressure: drop old deltas, send snapshot if client behind.

## Interfaces and boundaries
## Interfaces
- Produces WORLD_SNAPSHOT and WORLD_DELTA protobuf messages for ws_gateway sending.

## File-by-file walkthrough (expected / required)
## Expected files
- `grid.go` — cell mapping and buckets
- `watchers.go` — neighbor offsets, watcher list
- `diff.go` — snapshot/delta assembly + no-change suppression
- `resync.go` — snapshot resend when behind

## Gotchas / failure modes
## Gotchas
- AOI membership updates only on cell boundary crossings.
- Keep tick_seq strictly increasing.

## Acceptance criteria
## Done when
- WORLD_DELTA is not sent when no changes occur.
- Slow client gets snapshot resync, not unbounded queue growth.

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
- Only modify/create files under: `internal/aoi/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
