# proto

**Purpose:** Protocol schema (Protobuf) and generation rules.

## Canon inputs (authoritative)
- 03_PROTOCOL_CONTRACT — CONSOLIDATED
- PROJECT KERNEL v0.1 — CONSOLIDATED

## Constraints
- Language: Go
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Only additive changes are allowed once IDs/tags are locked.
- `msg_type` numbers and enum IDs are canon; never renumber.
- Protobuf payloads must match the fixed WSS framing header.

## Interfaces and boundaries
## Interfaces and consumers
- Server: `internal/protocol/pb` (generated), `internal/router`, `internal/ws_gateway`
- Client: TS protobuf bindings (chosen in Decision Ledger)

## File-by-file walkthrough (expected / required)
## Expected files
- `proto/game.proto` — canonical message + enum schema (IDs/tags stable)
- `proto/README.md` — generation instructions and output paths


## Gotchas / failure modes
## Common failure modes
- Renumbering fields breaks backwards compatibility and determinism.
- Divergence between server pb and client pb causes silent decode failures.

## Acceptance criteria
## Done when
- `protoc` generates code successfully.
- Frame decode/encode uses these messages consistently across server/client.

### Prompt seed for this subdirectory (for later)
Use this as the nucleus for a per-subdir generator prompt.

**Required attachments**
- `docs/CANON_LOCK.md`
- `docs/DECISION_LEDGER.md`
- `docs/STATE_HANDOFF.md` (latest)
- `TREE.txt`
- Any existing files under this subdir
- The governing design docs for this subdir (see “Canon inputs” above)

**Scope lock**
- Only modify/create files under: `proto/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
