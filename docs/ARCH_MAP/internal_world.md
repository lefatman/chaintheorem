---
owner: internal/world
status: IN_PROGRESS
generated_files:
  - internal/world/types.go
  - internal/world/store.go
  - internal/world/intents.go
  - internal/world/tick.go
touchpoints:
  - docs/ARCH_MAP/README.md
  - docs/STATE_HANDOFF.md
last_updated: 2026-02-02
---

# internal/world

**Purpose:** 10 Hz overworld simulation + entity store + movement intents.

## Canon inputs (authoritative)
- 04_AOI_AND_REPLICATION — CONSOLIDATED
- MVP System Architecture — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- 10 Hz tick; apply last movement intent per player.
- Stable entity_id allocation; deterministic ordering for diff production.
- Entity store should avoid per-tick allocations; prefer reusable slices.

## Interfaces and boundaries
## Interfaces
- Consumed by `internal/aoi` to compute watcher membership and diffs.
- Receives intents via router handler for WORLD_MOVE_INTENT.

## File-by-file walkthrough (expected / required)
## Expected files
- `types.go` — entity structs, kinds
- `store.go` — entity registry (stable IDs)
- `intents.go` — per-player last intent
- `tick.go` — 10 Hz loop / update step

## Gotchas / failure modes
## Gotchas
- Don’t scan all entities per player; AOI handles that.
- Avoid floats; use tile ints.

## Acceptance criteria
## Done when
- With two clients connected, movement produces correct authoritative positions at 10 Hz.

## Generated/Modified Files
- `internal/world/types.go`
- `internal/world/store.go`
- `internal/world/intents.go`
- `internal/world/tick.go`

## Interfaces / Contracts
- `EntitySource` for AOI reads (`AppendEntities`, `EntityByID`).
- `MoveIntentSink` for router/world input (`SetMoveIntent`).
- `Store` with deterministic entity storage + ID allocation.
- `IntentStore` with stable, sorted per-player intents.

## Algorithmic Invariants Implemented
- Stable monotonic entity IDs; deterministic iteration order.
- Movement applies last stored intent per player each tick.
- Intent application avoids per-tick allocations and map lookups in hot loops.

## Remaining Work
- Wire router/world handler to decode move intents and bind player entities.
- Integrate AOI reads and delta replication using the store.

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
- Only modify/create files under: `internal/world/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
