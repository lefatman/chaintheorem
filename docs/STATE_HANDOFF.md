# STATE HANDOFF — Batch 01 (Proto + Framing)

## What this batch created / updated (scope-locked)
### Protobuf sources
- `proto/game.proto`
  - Implements the protocol-contract message schemas (Hello/Welcome, overworld, chat, battle).
  - Adds minimal schemas for msg types present in the contract but not defined in the sketch:
    - `Ping{}` / `Pong{}`
    - `Error{ code, text }` (ledgered)
  - Canon enums included with the exact contract IDs:
    - ElementId (0..4)
    - ItemId (1..7) plus required proto3 zero value
    - AbilityId (1..8) plus required proto3 zero value
    - Dir4 (0..3), BattleActionType (0..1), TimelineEventType (0..7)

### Wire framing
- `internal/net/frame/frame.go`
  - Little-endian framing:
    - u16 msg_type
    - u32 payload_len
    - payload bytes (protobuf)
  - Bounds checks + max payload guard (default aligns to `config/server.json` ws.read_limit_bytes).

### Protocol constants for server code
- `internal/protocol/msgtypes.go`
  - Canon msg_type constants (u16) exactly matching the contract table.
- `internal/protocol/enums.go`
  - Canon IDs for ElementId/AbilityId/ItemId.
  - TimelineEventType constants.
  - **PieceType numeric IDs were unspecified**; internal mapping is ledgered (DECISION 0006).

### Decision ledger updates (append-only)
- Added:
  - DECISION 0005 (PING/PONG/ERROR protobuf shapes)
  - DECISION 0006 (PieceType numeric IDs)

## What this batch did NOT do (by design)
- No generated protobuf Go code yet (run `make proto` after installing protoc + protoc-gen-go).
- No WS gateway / router / handlers yet (Batch 02).
- No changes to protocol msg_type IDs, AOI tick, battle determinism, or capture-only rules.

## Verification steps (manual)
1. Inspect `proto/game.proto` for enum values and msg schemas.
2. Run protobuf generation:
   - macOS/Linux: `make proto`
   - Windows: `powershell -ExecutionPolicy Bypass -File .\scripts\gen_proto.ps1`
3. Sanity-check frame encode/decode:
   - HeaderLen=6, little-endian, strict length match.

## Next batch (Batch 02) must do
1. Add WS gateway skeleton (coder/websocket) that:
   - Reads full binary frames from WS
   - Uses `frame.Decode` to parse header + payload
   - Routes by `protocol.MsgType`
   - Applies per-conn write queue backpressure policy from config
2. Implement a minimal message router package layout under `internal/` (world/chat/battle) without inventing new protocol.
3. Implement config loaders (server.json + gameplay.json) with fail-fast validation and stable IDs.
