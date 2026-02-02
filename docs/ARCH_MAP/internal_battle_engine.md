# internal/battle_engine

**Purpose:** Deterministic chess legality + elements/abilities/items + timeline generation.

## Canon inputs (authoritative)
- Battle Engine Specification — CONSOLIDATED
- 07_GAMEPLAY_SYSTEMS — CONSOLIDATED
- 03_PROTOCOL_CONTRACT — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Deterministic baseline chess legality (including check) unless overridden.
- One-ply resolution order (fixed):
  1) validate input
  2) primary action: MOVE or CHAIN_KILL
  3) post-move Block Path set (if applicable)
  4) offensive post-capture triggers (Double/Quantum/Necromancer) with element ordering
  5) Poisoned Dagger reaction
  6) terminal state computation
- Deterministic RNG:
  - Lightning misfire vs Air/Wind
  - Quantum victim selection
  Seed and outcomes must be explicit in timeline events.
- Redo: rewind exactly 2 plies with snapshot restore and EV_REDO_REWIND event.

## Interfaces and boundaries
## Interfaces
- Consumed by `internal/battle_mgr` as a pure deterministic engine.
- Emits protobuf-compatible TimelineEvent sequences (or internal equivalent that maps 1:1).

## File-by-file walkthrough (expected / required)
## Expected files (suggested split)
- `ids.go` — ranks, piece types, constants
- `rng.go` — deterministic PRNG
- `board.go` — board representation + encode/decode snapshots
- `state.go` — per-piece and side-level state (charges, blocked dirs)
- `timeline.go` — timeline builder + event emit helpers
- `history.go` — ring buffer of ply-start snapshots (>= 2)
- `legality.go` — chess legality
- `resolve_move.go` — MOVE action resolution
- `resolve_chain_kill.go` — CHAIN_KILL resolution + Earth nullification
- `abilities_defensive.go` — Block Path/Stalwart/Belligerent/Redo
- `abilities_offensive.go` — Double/Quantum/Necromancer
- `items.go` — Poisoned Dagger + Solar Necklace usage
- `elements.go` — passives and interaction rules
- `apply_turn.go` — single entrypoint ApplyTurn(input)->timeline+new_state

## Gotchas / failure modes
## Gotchas
- NEVER introduce HP/damage/heal semantics.
- Don’t iterate over maps for victim selection; use deterministic ordering.
- “Fire vs Water” fizzles offensive abilities but still allows base captures.

## Acceptance criteria
## Done when
- Same battle seed + same inputs produce identical timelines.
- Redo correctly rewinds 2 plies and requires resubmission of the rewound turn.

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
- Only modify/create files under: `internal/battle_engine/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
