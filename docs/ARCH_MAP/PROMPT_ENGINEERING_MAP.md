# Per-Subdirectory Prompt Engineering Map

This document explains how to generate *future* per-subdir “code batches” safely.

## Core pattern
Every per-subdir prompt should:
 Require updating docs/STATE_HANDOFF.md.

## Minimal required updates for any batch
- docs/CANON_LOCK.md
- docs/DECISION_LEDGER.md
- docs/STATE_HANDOFF.md (latest)
- TREE.txt
- The governing canon docs for the subdir
- Any existing files under the subdir prefix you will modify


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
