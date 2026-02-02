---
status: done
owner: proto
generated_files:
  - proto/game.proto
  - proto/README.md
touchpoints:
  - internal/protocol/msgtypes.go
  - internal/protocol/enums.go
depends_on:
  - 00_global_contract
last_updated: 2026-02-02
---

# proto

**Purpose:** Protocol schema (Protobuf) and generation rules.

## What exists now (file-by-file)
- `proto/game.proto`
  - Canon msg types and message payloads.
  - Includes Ping/Pong empty messages and Error{code,text} (DECISION 0005).
  - Enums are canon-aligned with protocol IDs.
- `proto/README.md`
  - Protobuf generation instructions and output locations.

## Interfaces / exports
- Defines wire schemas for all msg types in `internal/protocol/msgtypes.go`.
- Generated Go types (when `make proto` is run) live under `internal/proto/gen` per `proto/README.md`.

## Constraints / invariants
- IDs and enum values are canonical and must never be renumbered.
- Wire framing is fixed-size header + protobuf payload (see `internal/net/frame`).

## Remaining work
- None in this module.

