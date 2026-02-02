# internal/protocol

**Purpose:** Shared constants/enums/msg_type registry and pb integration boundary.

## Canon inputs (authoritative)
- 03_PROTOCOL_CONTRACT — CONSOLIDATED
- PROJECT KERNEL v0.1 — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Centralized numeric constants: msg_type, ElementId/AbilityId/ItemId, piece ranks, event types.
- Keep generated pb types isolated (e.g., `internal/protocol/pb`).

## Interfaces and boundaries
## Interfaces
- Imported by world/aoi/chat/battle modules and router; must not import those modules back (no cycles).

## File-by-file walkthrough (expected / required)
## Expected files
- `internal/protocol/msgtypes.go` — msg_type constants
- `internal/protocol/enums.go` — Element/Ability/Item/PieceType/etc constants
- `internal/protocol/pb/` — generated protobuf output (or a README if generation is externalized)

## Gotchas / failure modes
## Gotchas
- Any mismatch vs proto schema creates runtime decode bugs.
- Keep “canonical ID tables” duplicated nowhere else.

## Acceptance criteria
## Done when
- All modules compile importing protocol constants.
- IDs match the canon tables in 03_PROTOCOL_CONTRACT / Gameplay Systems.

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
- Only modify/create files under: `internal/protocol/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
