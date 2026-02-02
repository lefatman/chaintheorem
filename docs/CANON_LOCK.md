# CANON LOCK (Batch 00)

This file restates **hard non-negotiables**. If a future batch conflicts with this, the future batch is wrong.

## Non-negotiables (hard)
1. **Authoritative Go server.** Clients are thin renderers/UI only.
2. **Transport**
   - HTTPS: auth endpoints + static assets
   - WSS: all realtime game traffic
3. **Binary protocol**
   - Fixed framing + Protobuf payloads
   - **No JSON on the hot path**
4. **Overworld**
   - **10 Hz** server tick
   - **Grid AOI**
   - **Delta diffs**
   - **No-change suppression**
5. **Battles**
   - Instanced **deterministic lockstep**
   - Clients send inputs
   - Server returns ordered **Outcome Timeline** events (**the only animation truth**)
6. **Capture-only semantics**
   - Pieces are either on-board or captured
   - **No HP / damage / healing** anywhere
7. **Efficiency guardrails**
   - Minimal allocations in hot loops
   - Bounded queues (no unbounded per-client buffers)
   - Stable IDs
   - Deterministic ordering
8. **Deterministic randomness**
   - If RNG exists, it is seeded at battle start
   - All RNG outcomes must be visible via outcome events

## Ownership reminder (cross-doc)
- Protocol contract owns wire shapes + msg types.
- AOI doc owns AOI logic + delta rules.
- Battle engine spec owns battle legality + deterministic resolution.
- Gameplay systems owns element/ability/item catalogs + IDs + loadout rules.
- Accounts/persistence doc owns auth + session + DB schema.
