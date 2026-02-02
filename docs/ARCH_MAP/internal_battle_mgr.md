# internal/battle_mgr

**Purpose:** Battle instance lifecycle, idempotency, reconnect, routing inputs to engine.

## Canon inputs (authoritative)
- Battle Engine Specification — CONSOLIDATED
- 07_GAMEPLAY_SYSTEMS — CONSOLIDATED
- 03_PROTOCOL_CONTRACT — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Owns instance lifecycle; engine is embedded but treated as deterministic core.
- Idempotency: repeated input for same (battle_id, expected turn_seq) returns same timeline or same rejection within same generation.
- Reconnect: can resend last timeline + expected turn_seq; snapshots after rewinds.

## Interfaces and boundaries
## Interfaces
- Router: BATTLE_TURN_INPUT -> battle_mgr
- Persist: loadouts/progression/xp awards
- ws_gateway: send reliability (timelines never dropped)

## File-by-file walkthrough (expected / required)
## Expected files
- `manager.go` — battle registry + lifecycle
- `instance.go` — per-instance state (expected turn_seq, seed, players, generation)
- `start.go` — start instance; emit BATTLE_START
- `handle_input.go` — validate turn_seq; call engine; emit BATTLE_OUTCOME_TIMELINE
- `end.go` — end conditions; emit BATTLE_END; award XP (Pot of Hunger multiplier)

## Gotchas / failure modes
## Gotchas
- After EV_REDO_REWIND, expected turn_seq becomes (timeline.turn_seq - 1) per protocol.
- Ensure old timelines for rewound turns are superseded (generation increments).

## Acceptance criteria
## Done when
- Two clients can play a battle end-to-end with timelines driving the client visuals.

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
- Only modify/create files under: `internal/battle_mgr/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
