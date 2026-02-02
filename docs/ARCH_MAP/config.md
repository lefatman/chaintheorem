# config

**Purpose:** External tunables: server.json, gameplay.json.

## Canon inputs (authoritative)
- PROJECT KERNEL v0.1 — CONSOLIDATED
- 07_GAMEPLAY_SYSTEMS — CONSOLIDATED
- 04_AOI_AND_REPLICATION — CONSOLIDATED

## Constraints
- Language: JSON
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- All tunables must live here; code reads configs at boot.
- gameplay.json must encode the canonical IDs and constraints (no hidden constants).
- server.json holds tick rates and AOI parameters.

## Interfaces and boundaries
## Interfaces
- Loaded by cmd/server and passed into world/aoi/battle/loadout validation.

## File-by-file walkthrough (expected / required)
## Expected files
- `config/server.json`
- `config/gameplay.json`
- (optional) schema docs: `config/README.md`

## Gotchas / failure modes
## Gotchas
- Never let gameplay IDs drift between config and protocol constants.

## Acceptance criteria
## Done when
- Changing config values changes behavior without code edits (within canon bounds).

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
- Only modify/create files under: `config/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
