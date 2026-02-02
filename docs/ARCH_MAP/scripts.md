# scripts

**Purpose:** Proto generation and dev utilities.

## Canon inputs (authoritative)
- PROJECT KERNEL v0.1 — CONSOLIDATED

## Constraints
- Language: Shell/PowerShell
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Scripts are helpers; never required for runtime.
- Proto generation must be deterministic and path-stable.

## Interfaces and boundaries
## Interfaces
- Used by developers only; must not be called by server at runtime.

## File-by-file walkthrough (expected / required)
## Expected files
- `scripts/gen_proto.sh`
- `scripts/gen_proto.ps1`
- (optional) `scripts/dev.sh`, `scripts/dev.ps1`

## Gotchas / failure modes
## Gotchas
- Keep Windows and *nix parity where possible.

## Acceptance criteria
## Done when
- `gen_proto` produces server/client bindings reproducibly.

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
- Only modify/create files under: `scripts/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
