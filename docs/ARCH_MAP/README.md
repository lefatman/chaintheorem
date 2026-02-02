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
