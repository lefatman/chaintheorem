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
