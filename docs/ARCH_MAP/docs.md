# docs

**Purpose:** Canon lock, decision ledger, state handoffs, and this architecture map.

## Canon inputs (authoritative)
- PROJECT KERNEL v0.1 — CONSOLIDATED

## Constraints
- Language: Markdown/Docs
- Must obey [Global Contracts](./00_global_contract.md) and the Canon Ownership Map.
- No silent changes; any necessary decision is appended to `docs/DECISION_LEDGER.md`.

## Algorithms and invariants
- Docs enforce canon and prevent drift.
- Decision Ledger is append-only.
- State handoff summarizes each batch output.

## Interfaces and boundaries
## Interfaces
- All build prompts treat these docs as binding inputs.

## File-by-file walkthrough (expected / required)
## Expected files
- `docs/CANON_LOCK.md`
- `docs/DECISION_LEDGER.md`
- `docs/STATE_HANDOFF.md`
- `docs/ARCH_MAP/*` (this pack)

## Gotchas / failure modes
## Gotchas
- Never silently rewrite Decision Ledger history.

## Acceptance criteria
## Done when
- Any non-canon choice is captured as an explicit ledger entry.

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
- Only modify/create files under: `docs/`  
- Append-only: `docs/DECISION_LEDGER.md`

**Hard rules**
- No silent changes: if not specified, add a Decision Ledger entry.
- No truncation: never output partial files; defer files if needed.
- Output full files only, each prefixed with `// File: path`.

**Task**
- Implement/extend the subdir exactly as described in this document, and update `docs/STATE_HANDOFF.md`.
