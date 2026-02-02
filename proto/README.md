# Protobufs (Batch 00)

This batch does **not** include `.proto` definitions yet and does **not** generate code.
It only establishes where protobuf sources will live and where generated output must go.

## Source location
Place protocol `.proto` files under:
- `proto/` (e.g., `proto/game.proto`, `proto/world.proto`, `proto/battle.proto`)

## Generated output location (Go)
Generation target (fixed):
- `internal/proto/gen/`

Generation command:
- `./scripts/gen_proto.sh` (macOS/Linux)
- `./scripts/gen_proto.ps1` (Windows PowerShell)

## Requirements
- `protoc` must be installed
- `protoc-gen-go` must be installed and on PATH:
  - `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

## Notes
- Message IDs and enums must match the **Protocol Contract — CONSOLIDATED**.
- Wire framing is **not** protobuf-delimited; framing is a custom header + protobuf payload per message type.
