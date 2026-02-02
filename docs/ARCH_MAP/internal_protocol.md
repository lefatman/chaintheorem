---
status: done
owner: internal/protocol
generated_files:
  - internal/protocol/msgtypes.go
  - internal/protocol/enums.go
touchpoints:
  - internal/router/handlers.go
  - internal/ws_gateway/conn.go
depends_on:
  - 00_global_contract
  - proto
last_updated: 2026-02-02
---

# internal/protocol

**Purpose:** Shared constants/enums/msg_type registry and pb integration boundary.

## What exists now (file-by-file)
- `internal/protocol/msgtypes.go`
  - Canon msg_type constants used in frame headers.
- `internal/protocol/enums.go`
  - Canon ElementId/AbilityId/ItemId constants.
  - PieceType numeric IDs (DECISION 0006).
  - Direction, action, and timeline enums.

## Interfaces / exports
- `protocol.MsgType` constants for routing and framing.
- Typed enums for canonical IDs used by gameplay code.

## Constraints / invariants
- Numeric IDs must match `proto/game.proto` and contract tables.
- No imports from downstream modules to avoid cycles.

## Remaining work
- None in this module.

