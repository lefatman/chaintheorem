---
status: done
owner: internal/ws_gateway
generated_files:
  - internal/ws_gateway/config.go
  - internal/ws_gateway/server.go
  - internal/ws_gateway/conn.go
  - internal/ws_gateway/queue.go
  - internal/ws_gateway/errors.go
touchpoints:
  - internal/router/router.go
  - internal/router/handlers.go
  - internal/net/frame/frame.go
  - internal/protocol/msgtypes.go
  - cmd/server/main.go
depends_on:
  - internal_net_frame
  - internal_protocol
  - internal_router
last_updated: 2026-02-02
---

# internal/ws_gateway

**Purpose:** WSS lifecycle, read/write loops, backpressure policy.

## What exists now (file-by-file)
- `config.go`
  - Gateway runtime config (read limits, queue size, timeouts).
- `server.go`
  - HTTP handler that upgrades to WebSocket and starts connection loops.
- `conn.go`
  - Read loop parses frames and dispatches by msg_type.
  - Write loop drains a bounded queue with per-frame timeouts.
  - Ping/Pong handled in-gateway for keepalive.
- `queue.go`
  - Bounded ring buffer plus a single droppable slot for coalesced deltas.
- `errors.go`
  - Sentinel errors for backpressure and lifecycle failures.

## Interfaces / exports
- `Server` implements `http.Handler` for the WS endpoint.
- `Config` defines runtime tuning parameters.

## Backpressure policy (implemented)
- Droppable: `MSG_WORLD_DELTA` (coalesced into a single pending slot when enabled).
- Non-droppable: all other outbound msg types.
- If non-droppable frames cannot be queued, the connection is closed with policy violation.

## Constraints / invariants
- Read loop only accepts binary messages and validates framing.
- No protobuf parsing in the gateway; payloads are forwarded raw.

## Remaining work
- None in this module.

