---
status: done
owner: internal/router
generated_files:
  - internal/router/router.go
  - internal/router/handlers.go
touchpoints:
  - internal/ws_gateway/conn.go
  - internal/app/app.go
depends_on:
  - internal_protocol
last_updated: 2026-02-02
---

# internal/router

**Purpose:** Dispatch frames to module handlers; schema-level validation only.

## What exists now (file-by-file)
- `router.go`
  - Fixed-size dispatch table keyed by `protocol.MsgType`.
  - Emits sentinel errors for unhandled or unauthenticated messages.
- `handlers.go`
  - Minimal handler interfaces for auth/world/chat/battle.
  - Registration helpers per module.

## Interfaces / exports
- `Router` with `Register` and `Dispatch`.
- Module handler interfaces: `AuthHandler`, `WorldHandler`, `ChatHandler`, `BattleHandler`.

## Constraints / invariants
- Dispatch uses an array table (no map iteration in hot path).
- Payloads stay as `[]byte` (no protobuf dependency here).

## Remaining work
- None in this module.

