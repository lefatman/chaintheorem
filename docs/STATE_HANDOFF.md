# STATE HANDOFF — Batch 03 (WS Gateway Deadlines)

## What changed (scope-locked)
- `internal/ws_gateway/conn.go`
  - Replaced per-iteration context timeouts with websocket read/write deadlines.
- `docs/ARCH_MAP/internal_ws_gateway.md`
  - Updated generated file list and invariants for deadline-based IO.
- `docs/ARCH_MAP/README.md`
  - Added last-updated timestamp.

## Decisions appended
- None.

## How to validate
1. `gofmt -w .`
2. `go test ./...`

## Next module to work on
- `internal_persist` (docs/ARCH_MAP/internal_persist.md)
