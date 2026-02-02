# ARCH_MAP_FULL


---

<!-- FILE: 00_global_contract.md -->

# Global Contracts and Invariants

This document is the “gravity” that prevents drift.

## Wire + simulation invariants
- **Realtime is WSS binary frames** carrying **Protobuf payloads**.
- Frame header is **little-endian**:
  - `u16 msg_type`
  - `u32 payload_len`
  - `payload_len` bytes payload
- **No JSON on hot path** (realtime). JSON is acceptable only on HTTPS endpoints.
- **Overworld**: authoritative **10 Hz** tick; AOI is **grid-based**; deltas are **suppressed if no change**.
- **Battles**: instanced **deterministic lockstep**:
  - clients send **inputs**
  - server returns **Outcome Timeline** events (only animation truth)
- **Capture-only semantics**: no HP/damage/heal anywhere.

## Determinism rules
- Any randomness is deterministic PRNG seeded at battle start.
- Any RNG outcome is emitted as an explicit timeline event.
- Deterministic ordering: do not iterate over Go maps for anything that affects output.

## “No silent changes” protocol
If you must decide something that is not specified by the canon docs, you must:
1) choose the smallest / simplest option,
2) append a new Decision Ledger entry describing the choice and rationale,
3) ensure downstream modules consume the same decision.

## msg_type registry (canonical)
| msg_type | Direction | Name | Purpose |
| --- | --- | --- | --- |
| 1 | C->S | HELLO | Authenticate WSS connection using session token. |
| 2 | S->C | WELCOME | Bind connection to player_id; initial routing info. |
| 3 | C->S | PING | Keepalive. |
| 4 | S->C | PONG | Keepalive response. |
| 10 | C->S | WORLD_MOVE_INTENT | Request overworld movement (authoritative on server tick). |
| 11 | S->C | WORLD_SNAPSHOT | Full AOI state (initial or resync). |
| 12 | S->C | WORLD_DELTA | AOI delta diff at 10 Hz (no-change suppressed). |
| 20 | C->S | CHAT_SEND | Send chat message. |
| 21 | S->C | CHAT_EVENT | Chat broadcast. |
| 30 | S->C | BATTLE_START | Enter battle instance; includes seed and initial board. |
| 31 | C->S | BATTLE_TURN_INPUT | Submit one ply input for deterministic lockstep. |
| 32 | S->C | BATTLE_OUTCOME_TIMELINE | Authoritative ordered events for one ply. |
| 33 | S->C | BATTLE_END | Terminal battle result + rewards summary. |
| 250 | S->C | ERROR | Error / rejection. |


## Elements, Abilities, Items (canonical tables)
### Elements
| ElementId | Name | Passive rules (canonical) |
| --- | --- | --- |
| 0 | Water | Consumable ability counters are doubled. This doubling is negated when fighting a Lightning army. |
| 1 | Fire | Offensive abilities resolve first. Fire offensive abilities are ineffective against Water armies. |
| 2 | Earth | Remote offensive capture abilities are nullified. This nullification is negated by Fire armies. |
| 3 | Air/Wind | Negates defensive abilities and can move over pieces. Air/Wind passives are negated by Earth armies. |
| 4 | Lightning | Army abilities are slottable at the piece-type level (no Multitasker’s Schedule needed). Abilities have a 50% chance to misfire against Air/Wind armies. |


### Items
| ItemId | Item | Slot cost | Effect |
| --- | --- | --- | --- |
| 1 | Multitasker’s Schedule | 1 | Allows piece-type-level slotting of army abilities (non-Lightning armies). |
| 2 | Poisoned Dagger | 1 | When a piece is captured, removes the capturing piece if the capturer is of lower or equal rank. |
| 3 | Dual Adept’s Gloves | 1 | Adds +1 army-ability slot. |
| 4 | Triple Adept’s Gloves | 2 | Adds +2 army-ability slots; blocks Dual Adept’s Gloves and Headmaster Ring. |
| 5 | Headmaster Ring | 3 | Adds +3 army-ability slots; blocks Dual and Triple Adept’s Gloves. |
| 6 | Pot of Hunger | 1 | Doubles XP gained from winning a match. |
| 7 | Solar Necklace | 1 | Can top up a consumable ability up to 3 times per match. |


### Abilities
| AbilityId | Ability | Scope | Category | Consumable | Canonical effect |
| --- | --- | --- | --- | --- | --- |
| 1 | Block Path |army-wide | defensive | no | After moving, choose a cardinal direction; this piece cannot be captured from that direction until it moves again. |
| 2 | Stalwart | army-wide | defensive | no | Pieces of this type cannot be captured by a lower-rank capturer. |
| 3 | Belligerent | army-wide| defensive | no | Pieces of this type cannot be captured by a higher-rank capturer. |
| 4 | Redo | army-wide | defensive | yes (1 per piece; doubled for Water vs non-Lightning) | When a piece of this type is captured, rewind exactly 2 plies (the capturing ply and the defender ply immediately before it), restoring the full match state to the start of the defender's prior turn. The defender then replays that prior turn with a different move. All effects of those two plies are reversed. Spend 1 Redo charge on the captured piece; the spent charge remains spent after the rewind. |
| 5 | Double Kill | army-wide | offensive | no | On capture, remove one neighboring enemy piece of equal or lower rank; if none exists, no effect. |
| 6 | Quantum Kill | army-wide | offensive | no | On capture, remove one random enemy piece of equal or lower rank. |
| 7 | Chain Kill | army-wide | offensive | no | Active capture: a piece can piggyback on an adjacent allied piece to capture a target as if it were on the ally’s square (remote capture). |
| 8 | Necromancer | army-wide | offensive | yes (side pool; doubled for Water vs non-Lightning) | If your piece captures a higher-rank piece, and you have any captured pieces, return one eligible captured friendly piece to the board at its capture square. |



---

<!-- FILE: PROMPT_ENGINEERING_MAP.md -->

# Per-Subdirectory Prompt Engineering Map

This document explains how to generate *future* per-subdir “code batches” safely.

## Core pattern
Every per-subdir prompt should:
1) Declare itself STATELESS.
2) Declare attachments as the universe.
3) Enforce a SCOPE LOCK limited to the subdir prefix (plus append-only Decision Ledger).
4) Demand a BATCH MANIFEST and full-file outputs only.
5) Require updating docs/STATE_HANDOFF.md.

## Minimal required attachments for any batch
- docs/CANON_LOCK.md
- docs/DECISION_LEDGER.md
- docs/STATE_HANDOFF.md (latest)
- TREE.txt
- The governing canon docs for the subdir
- Any existing files under the subdir prefix you will modify

## “Stop if missing” guard
Your prompt must include:
- If a needed file is not attached, STOP and output only:
  MISSING INPUTS:
  - <path>

## Output hygiene
- Never emit diffs. Full files only.
- Never start a file you can’t finish (no truncation).
- Prefer small files (<200 lines) and more files over monoliths.

## Suggested batch slicing (server)
- proto + framing (already done)
- ws_gateway + router
- persist + migrations
- auth + httpapi
- world + aoi
- chat
- battle_engine foundations
- battle_engine rules
- battle_mgr lifecycle
- loadout validation

## Suggested batch slicing (client)
- boot + frame decode + WS
- overworld: state apply + interpolation + render
- battle: timeline player + UI
- PWA: service worker + caching + offline menu


---

<!-- FILE: README.md -->

# Architecture Map Pack (Stateless, Drift-Resistant)

This pack reorganizes the consolidated design docs into **per-subdirectory Markdown constraints** so you can later generate **per-subdir prompts** without the model inventing structure.

## Canon non-negotiables (global)
- Authoritative Go server; clients are thin renderers/UI only.
- Transport: HTTPS for assets/auth; WSS for realtime.
- Binary protocol: fixed framing + Protobuf; no JSON on hot path.
- Overworld: 10 Hz tick, grid AOI, delta diffs, no-change suppression.
- Battles: instanced deterministic lockstep; inputs in, Outcome Timeline out (only animation truth).
- Capture-only semantics: pieces are on-board or captured. No HP/damage/heal anywhere.
- Efficiency: minimal allocations in hot loops; bounded queues; stable IDs; deterministic ordering.
- Deterministic randomness only: seeded at battle start; outcomes disclosed via timeline events.

## Canon ownership map (who decides what)
| Doc | Owns | Must NOT own |
| --- | --- | --- |
| PROJECT KERNEL | Global defaults and non-negotiables | Detailed rules/protocol/DB schemas |
| MVP System Architecture | Runtime component boundaries and flows | Exact wire schemas or gameplay tables |
| PixiJS Client Plan | Rendering + UX + client responsibilities | Server rules/authority |
| 03_PROTOCOL_CONTRACT | Wire format + message definitions + sequencing | Battle legality algorithms |
| 04_AOI_AND_REPLICATION | AOI strategy + delta diff rules | Account schema details |
| Battle Engine Spec | Battle legality + deterministic resolution | Overworld AOI algorithms |
| Accounts/Persistence/Social | DB schema + auth/social flows | Element/ability item definitions |
| Gameplay Systems | Elements + abilities + items + loadout rules | Protocol framing |


## Runtime module ownership (server)
| Module | Owns | Does NOT own |
| --- | --- | --- |
| ws_gateway | WSS connection lifecycle, frame parse/encode, backpressure | Game rules, DB writes |
| router | Dispatch by message type; validation of basic schema invariants | Battle legality, AOI logic |
| auth | Register/login/reset flows; token issuance/verification | AOI, battle simulation |
| persist | SQL schema + queries + transactions | Business rules (except data constraints) |
| world | 10 Hz overworld sim, entity storage, AOI membership | Battle rules |
| aoi | Spatial grid, watchers, diff assembly | WebSocket framing |
| chat | Channel membership and broadcast policy | Auth storage |
| battle_mgr | Instance lifecycle, matchmaking hooks, routing inputs to instances | Overworld AOI |
| battle_engine | Deterministic chess rules + abilities/items/elements + timelines | Account registration |


## Subdirectory index
- [proto](./proto.md) — Protocol schema (Protobuf) and generation rules.
- [internal/net/frame](./internal_net_frame.md) — Binary frame codec (u16 msg_type + u32 payload_len LE) used by ws_gateway.
- [internal/protocol](./internal_protocol.md) — Shared constants/enums/msg_type registry and pb integration boundary.
- [internal/ws_gateway](./internal_ws_gateway.md) — WSS lifecycle, read/write loops, backpressure policy, session binding handshake.
- [internal/router](./internal_router.md) — Dispatch frames to module handlers; schema-level validation only.
- [internal/app](./internal_app.md) — Dependency wiring + server bootstrap composition.
- [internal/auth](./internal_auth.md) — Register/login/reset + session tokens, password hashing.
- [internal/httpapi](./internal_httpapi.md) — HTTPS JSON endpoints for auth + loadout editing + dev hooks.
- [internal/persist](./internal_persist.md) — DB schema, migrations, repositories; SQLite dev + Postgres prod.
- [internal/world](./internal_world.md) — 10 Hz overworld simulation + entity store + movement intents.
- [internal/aoi](./internal_aoi.md) — Grid AOI, watcher sets, diff assembly (snapshot/delta), resync.
- [internal/chat](./internal_chat.md) — Global chat broadcast + rate limiting.
- [internal/battle_engine](./internal_battle_engine.md) — Deterministic chess legality + elements/abilities/items + timeline generation.
- [internal/battle_mgr](./internal_battle_mgr.md) — Battle instance lifecycle, idempotency, reconnect, routing inputs to engine.
- [config](./config.md) — External tunables: server.json, gameplay.json.
- [scripts](./scripts.md) — Proto generation and dev utilities.
- [client](./client.md) — PixiJS PWA thin client: rendering, networking decode, interpolation, battle timeline animation.
- [docs](./docs.md) — Canon lock, decision ledger, state handoffs, and this architecture map.


---

<!-- FILE: client.md -->

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


---

<!-- FILE: config.md -->

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


---

<!-- FILE: docs.md -->

# docs

**Purpose:** Canon lock, decision ledger, state handoffs, and this architecture map.

## Canon inputs (authoritative)
- PROJECT KERNEL v0.1 — CONSOLIDATED

## Constraints
- Language: Markdown/Docs
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Docs enforce canon and prevent drift.
- Decision Ledger is append-only.
- State handoff summarizes each batch output.

## Interfaces and boundaries
## Interfaces
- All build prompts treat these docs as binding inputs.

## File-by-file walkthrough (expected / required)
## Expected files
- `docs/CANON_LOCK.md`
- `docs/DECISION_LEDGER.md`
- `docs/STATE_HANDOFF.md`
- `docs/ARCH_MAP/*` (this pack)

## Gotchas / failure modes
## Gotchas
- Never silently rewrite Decision Ledger history.

## Acceptance criteria
## Done when
- Any non-canon choice is captured as an explicit ledger entry.

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
- Only modify/create files under: `docs/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.


---

<!-- FILE: internal_aoi.md -->

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


---

<!-- FILE: internal_app.md -->

# internal/app

**Purpose:** Dependency wiring + server bootstrap composition.

## Canon inputs (authoritative)
- MVP System Architecture — CONSOLIDATED
- PROJECT KERNEL v0.1 — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Composition root only: wire modules and start servers.
- No business logic beyond configuration and lifecycle.

## Interfaces and boundaries
## Interfaces
- Called from `cmd/server/main.go`.

## File-by-file walkthrough (expected / required)
## Expected files
- `app.go` — wire dependencies and expose Start/Stop
- (optional) `wiring.go` — small helpers to keep app.go short

## Gotchas / failure modes
## Gotchas
- Don’t hide decisions here; put them in configs or ledger.

## Acceptance criteria
## Done when
- `cmd/server` can start the server with this app wiring.

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
- Only modify/create files under: `internal/app/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.


---

<!-- FILE: internal_auth.md -->

# internal/auth

**Purpose:** Register/login/reset + session tokens, password hashing.

## Canon inputs (authoritative)
- 06_ACCOUNTS_PERSISTENCE_SOCIAL — CONSOLIDATED
- PROJECT KERNEL v0.1 — CONSOLIDATED
- MVP System Architecture — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- HTTPS register/login/reset flows; WSS uses HELLO token binding.
- Password hashing via bcrypt or argon2id (decision is ledger-owned).
- Session tokens are opaque random bytes stored server-side with expiry (MVP).

## Interfaces and boundaries
## Interfaces
- Uses `internal/persist` for accounts/sessions.
- Exposed via `internal/httpapi` handlers.

## File-by-file walkthrough (expected / required)
## Expected files
- `service.go` — high-level auth operations
- `password.go` — hash + verify
- `tokens.go` — token generation + expiry
- (optional) `ratelimit.go` — basic anti-abuse

## Gotchas / failure modes
## Gotchas
- Never log raw passwords or session tokens.
- Token comparison should be constant-time if feasible.

## Acceptance criteria
## Done when
- Register/login works end-to-end; HELLO binds a WSS connection using issued token.

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
- Only modify/create files under: `internal/auth/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.


---

<!-- FILE: internal_battle_engine.md -->

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


---

<!-- FILE: internal_battle_mgr.md -->

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


---

<!-- FILE: internal_chat.md -->

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


---

<!-- FILE: internal_httpapi.md -->

# internal/httpapi

**Purpose:** HTTPS JSON endpoints for auth + loadout editing + dev hooks.

## Canon inputs (authoritative)
- 06_ACCOUNTS_PERSISTENCE_SOCIAL — CONSOLIDATED
- PROJECT KERNEL v0.1 — CONSOLIDATED
- MVP System Architecture — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- JSON over HTTPS is allowed (not hot path).
- Keep handlers thin: validate, call service, return response.

## Interfaces and boundaries
## Interfaces
- Uses `auth`, `persist`, `loadout`, `battle_mgr` (for dev hooks).

## File-by-file walkthrough (expected / required)
## Expected files
- `server.go` — mux/routes
- `auth_handlers.go` — register/login/reset
- `loadout_handlers.go` — GET/POST loadout (for MVP testing)
- (optional) `dev_handlers.go` — dev-only battle start hooks

## Gotchas / failure modes
## Gotchas
- Keep CORS and cookie/token policies consistent; document choices in ledger.

## Acceptance criteria
## Done when
- You can create an account and obtain a session token via HTTPS endpoints.

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
- Only modify/create files under: `internal/httpapi/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.


---

<!-- FILE: internal_net_frame.md -->

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


---

<!-- FILE: internal_persist.md -->

# internal/persist

**Purpose:** DB schema, migrations, repositories; SQLite dev + Postgres prod.

## Canon inputs (authoritative)
- 06_ACCOUNTS_PERSISTENCE_SOCIAL — CONSOLIDATED
- PROJECT KERNEL v0.1 — CONSOLIDATED
- MVP System Architecture — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- SQLite dev + Postgres prod; schema portability is required.
- Migration runner: version table + ordered migrations.
- Repos expose minimal CRUD; business rules stay outside persist.

## Interfaces and boundaries
## Canon schema (MVP extract)
```sql
accounts(user_id PK, email UNIQUE, username UNIQUE, pass_hash, created_at, last_login_at)
sessions(token PK, user_id, expires_at, created_at)
army_loadouts(user_id PK, element_id, army_ability_1..4, ability_*piece, item_1..4, updated_at)
progression(user_id PK, level, xp)
user_unlocks(user_id, flag_id, unlocked_at, PK(user_id, flag_id))
```

## File-by-file walkthrough (expected / required)
## Expected files
- `db.go` — open/health/ping
- `migrate.go` — migration runner
- `migrations/*.sql` — ordered schema files
- `accounts_repo.go`, `sessions_repo.go`, `loadouts_repo.go`, `progression_repo.go`, `unlocks_repo.go`

## Gotchas / failure modes
## Gotchas
- SQLite uses `INTEGER PRIMARY KEY` rowid; Postgres needs sequences/identity.
- Keep timestamps UTC integer seconds.

## Acceptance criteria
## Done when
- Fresh DB initializes via migrations on both SQLite and Postgres (or clearly separated dialect migrations).

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
- Only modify/create files under: `internal/persist/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.


---

<!-- FILE: internal_protocol.md -->

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


---

<!-- FILE: internal_router.md -->

# internal/router

**Purpose:** Dispatch frames to module handlers; schema-level validation only.

## Canon inputs (authoritative)
- 03_PROTOCOL_CONTRACT — CONSOLIDATED
- PROJECT KERNEL v0.1 — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Dispatch by msg_type.
- Validate schema invariants only (e.g., dx/dy in -1..1), not game legality.
- No heavy allocations; avoid map iteration for deterministic operations (use switch/array tables).

## Interfaces and boundaries
## Interfaces
- Upstream: ws_gateway passes (player_id, msg_type, payload)
- Downstream: module handlers: auth/world/aoi/chat/battle_mgr

## File-by-file walkthrough (expected / required)
## Expected files
- `router.go` — dispatch core
- `handlers.go` — module handler interfaces
- (optional) `errors.go` — ERROR message helpers

## Gotchas / failure modes
## Gotchas
- Router must not import heavy subsystems that cause cycles.
- Keep all msg_type handling in one obvious place (auditability).

## Acceptance criteria
## Done when
- Every msg_type in protocol has exactly one handler path or an explicit rejection.

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
- Only modify/create files under: `internal/router/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.


---

<!-- FILE: internal_world.md -->

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


---

<!-- FILE: internal_ws_gateway.md -->

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


---

<!-- FILE: proto.md -->

# proto

**Purpose:** Protocol schema (Protobuf) and generation rules.

## Canon inputs (authoritative)
- 03_PROTOCOL_CONTRACT — CONSOLIDATED
- PROJECT KERNEL v0.1 — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Only additive changes are allowed once IDs/tags are locked.
- `msg_type` numbers and enum IDs are canon; never renumber.
- Protobuf payloads must match the fixed WSS framing header.

## Interfaces and boundaries
## Interfaces and consumers
- Server: `internal/protocol/pb` (generated), `internal/router`, `internal/ws_gateway`
- Client: TS protobuf bindings (chosen in Decision Ledger)

## File-by-file walkthrough (expected / required)
## Expected files
- `proto/game.proto` — canonical message + enum schema (IDs/tags stable)
- `proto/README.md` — generation instructions and output paths


## Gotchas / failure modes
## Common failure modes
- Renumbering fields breaks backwards compatibility and determinism.
- Divergence between server pb and client pb causes silent decode failures.

## Acceptance criteria
## Done when
- `protoc` generates code successfully.
- Frame decode/encode uses these messages consistently across server/client.

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
- Only modify/create files under: `proto/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.


---

<!-- FILE: scripts.md -->

# scripts

**Purpose:** Proto generation and dev utilities.

## Canon inputs (authoritative)
- PROJECT KERNEL v0.1 — CONSOLIDATED

## Constraints
- Language: Shell/PowerShell
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Scripts are helpers; never required for runtime.
- Proto generation must be deterministic and path-stable.

## Interfaces and boundaries
## Interfaces
- Used by developers only; must not be called by server at runtime.

## File-by-file walkthrough (expected / required)
## Expected files
- `scripts/gen_proto.sh`
- `scripts/gen_proto.ps1`
- (optional) `scripts/dev.sh`, `scripts/dev.ps1`

## Gotchas / failure modes
## Gotchas
- Keep Windows and *nix parity where possible.

## Acceptance criteria
## Done when
- `gen_proto` produces server/client bindings reproducibly.

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
- Only modify/create files under: `scripts/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
