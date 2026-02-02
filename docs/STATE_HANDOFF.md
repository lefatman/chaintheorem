# STATE HANDOFF — Batch 06 (World Core)

## What this batch created / updated (scope-locked)
### World core scaffolding
- `internal/world/types.go`
- `internal/world/store.go`
- `internal/world/intents.go`
- `internal/world/tick.go`
  - Added deterministic entity store with stable IDs.
  - Added per-player last-move intent tracking and bindable player/entity indices.
  - Added tick step to apply intents without per-tick allocations.

### Documentation updates
- `docs/ARCH_MAP/internal_world.md`
- `docs/ARCH_MAP/README.md`
- `docs/STATE_HANDOFF.md`

## Decisions appended
- None.

## How to validate
1. `gofmt -w .`
2. `go test ./...`

## Next module to implement
- `docs/ARCH_MAP/internal_world.md` (finish wiring) or proceed to `docs/ARCH_MAP/internal_aoi.md` once world is DONE.
