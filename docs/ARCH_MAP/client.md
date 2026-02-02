# client

**Purpose:** PixiJS PWA thin client: rendering, networking decode, interpolation, battle timeline animation.

## Canon inputs (authoritative)
- PixiJS Client Architecture Plan — CONSOLIDATED
- 03_PROTOCOL_CONTRACT — CONSOLIDATED

## Constraints
- Language: TypeScript + Web
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Thin client: renderer/UI only.
- Decode WSS binary frames + protobuf messages.
- Overworld interpolation buffer 150–250ms; minimal or no prediction.
- Z-sort by (y, entity_id) for stable draw order.
- Battles: animate strictly from Outcome Timeline at 60 FPS; EV_REDO_REWIND uses snapshot to rewind visuals.

## Interfaces and boundaries
## Interfaces
- Protocol decode must match server msg_type and protobuf schema exactly.
- Client must not assume legality; server remains authoritative.

## File-by-file walkthrough (expected / required)
## Expected directory layout (suggested)
- `client/index.html`
- `client/manifest.json`
- `client/sw.js`
- `client/src/net/` (frame decode + WS client)
- `client/src/world/` (snapshot/delta apply + interpolation)
- `client/src/render/` (Pixi init, camera, sprite pooling, z-sort)
- `client/src/battle/` (board decode, timeline player, battle UI)
- `client/src/main.ts`

## Gotchas / failure modes
## Gotchas
- Avoid texture swaps; use atlases.
- Never reorder or infer missing timeline events.

## Acceptance criteria
## Done when
- Overworld shows moving entities smoothly.
- Battle timelines animate deterministically and rewinds resync correctly.

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
- Only modify/create files under: `client/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
