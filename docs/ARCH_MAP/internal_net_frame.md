---
status: done
owner: internal/net/frame
generated_files:
  - internal/net/frame/frame.go
touchpoints:
  - internal/ws_gateway/conn.go
depends_on:
  - 00_global_contract
  - proto
last_updated: 2026-02-02
---

# internal/net/frame

**Purpose:** Binary frame codec (u16 msg_type + u32 payload_len LE) used by ws_gateway.

## What exists now (file-by-file)
- `internal/net/frame/frame.go`
  - Encode/Decode helpers for the fixed header and payload.
  - Enforces payload max and strict length match.

## Interfaces / exports
- `Encode(dst, msgType, payload, maxPayload)`
- `Decode(buf, maxPayload)`
- `HeaderLen` and `DefaultMaxPayloadBytes`

## Constraints / invariants
- Little-endian header encoding.
- Strict bounds checks to prevent oversized payloads.

## Remaining work
- None in this module.

