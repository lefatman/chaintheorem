---
status: done
owner: internal/app
generated_files:
  - internal/app/app.go
touchpoints:
  - internal/ws_gateway/server.go
  - internal/router/router.go
  - cmd/server/main.go
  - internal/config/config.go
depends_on:
  - internal_router
  - internal_ws_gateway
  - internal_config
last_updated: 2026-02-02
---

# internal/app

**Purpose:** Dependency wiring + server bootstrap composition.

## What exists now (file-by-file)
- `app.go`
  - Constructs router and gateway.
  - Registers stub handlers that accept HELLO and disconnect for unimplemented modules.

## Interfaces / exports
- `New(serverCfg)` returns `*App` with `Router` and `Gateway`.

## Constraints / invariants
- No gameplay logic; composition only.
- Handlers should be safe no-ops or disconnect (no protobuf required yet).

## Remaining work
- Replace stub handlers with real auth/world/chat/battle modules.

