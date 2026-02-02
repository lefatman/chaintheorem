# STATE HANDOFF — Batch 02 (WS Gateway + Router)

## What this batch created / updated (scope-locked)
### WebSocket gateway
- `internal/ws_gateway/config.go`
- `internal/ws_gateway/server.go`
- `internal/ws_gateway/conn.go`
- `internal/ws_gateway/queue.go`
- `internal/ws_gateway/errors.go`
  - Accepts WebSocket connections (coder/websocket).
  - Reads binary frames, parses via `internal/net/frame`.
  - Dispatches by `protocol.MsgType` to the router.
  - Maintains bounded outbound queue with droppable coalescing for overworld deltas.

### Router
- `internal/router/router.go`
- `internal/router/handlers.go`
  - Central dispatch table keyed by msg_type.
  - Minimal handler interfaces for auth/world/chat/battle.

### App wiring
- `internal/app/app.go`
  - Composes router + gateway and registers stub handlers.
  - HELLO is accepted; other modules disconnect with unimplemented errors.

### Config loaders
- `internal/config/config.go`
- `internal/config/gameplay.go`
  - Strict JSON decoding with validation of required fields and canonical IDs.

### Server entrypoint
- `cmd/server/main.go`
  - Loads configs, starts HTTP `/health` and WS endpoint.
  - WS TLS enabled via `config/server.json`.

### Documentation updates
- `docs/ARCH_MAP/*` updated to reflect reality.
- `docs/ARCH_MAP/README.md` generation order and statuses.

## Backpressure policy (implemented)
- Droppable/coalescible: `MSG_WORLD_DELTA`.
- Non-droppable: all other outbound message types.
- If a non-droppable frame cannot be queued, the connection is closed with policy violation.

## Routing behavior (implemented)
- WS gateway decodes frames, routes by msg_type.
- Router dispatches to module handlers with raw payload bytes.
- No protobuf parsing in the gateway/router at this stage.

## What was NOT done (by design)
- No auth/session validation beyond accepting HELLO.
- No WELCOME/Error protobuf emission (protobuf generation not required yet).
- No world/chat/battle gameplay logic.
- No persistence/auth/httpapi modules.

## Verification steps (manual)
1. `gofmt -w internal/config/*.go internal/router/*.go internal/ws_gateway/*.go internal/app/*.go cmd/server/*.go`
2. `go test ./...` (will be minimal until more modules exist)
3. Run server with configs:
   - `go run ./cmd/server`
   - `curl http://localhost:8080/health` returns `ok`.

## Next batch requirements
- Implement persistence module (`internal/persist`) and wire into app/router.
- Implement auth module (`internal/auth`) and connect HELLO -> WELCOME flow.
- Update `docs/ARCH_MAP/internal_persist.md` and `docs/ARCH_MAP/internal_auth.md` to done.

**Next module doc to implement:** `docs/ARCH_MAP/internal_persist.md`

