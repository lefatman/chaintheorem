---
status: done
owner: internal/config
generated_files:
  - internal/config/config.go
  - internal/config/gameplay.go
touchpoints:
  - cmd/server/main.go
  - config/server.json
  - config/gameplay.json
depends_on:
  - 00_global_contract
last_updated: 2026-02-02
---

# internal/config

**Purpose:** Load and validate server/gameplay configuration at boot.

## What exists now (file-by-file)
- `config.go`
  - Typed `ServerConfig` with strict JSON decoding and validation.
- `gameplay.go`
  - Typed `GameplayConfig` with strict JSON decoding and canonical ID validation.

## Interfaces / exports
- `LoadServerConfig(path)`
- `LoadGameplayConfig(path)`

## Constraints / invariants
- Unknown JSON fields are rejected (fail-fast).
- Canonical IDs must be complete and sequential per config counts.

## Remaining work
- None in this module.

