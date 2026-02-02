# AGENTS.md — Codex Project Instructions (Battle Chess MMO)

This repository is built in **stateless, drift-resistant batches** guided by `docs/ARCH_MAP/*`.
Codex should treat these rules as **binding** for every run.

---

## 1) Canon hierarchy (source of truth)

**Highest priority**
1. `docs/CANON_LOCK.md` — non‑negotiables / invariants.
2. `proto/game.proto` + `internal/protocol/*` — **wire IDs and enums are binding**.
3. `docs/DECISION_LEDGER.md` — append‑only; any non‑canon choice must be recorded here.
4. `docs/STATE_HANDOFF.md` — latest work summary + next module.

**Architecture map**
- `docs/ARCH_MAP/00_global_contract.md` — project-wide constraints and decision mirroring.
- `docs/ARCH_MAP/README.md` — **explicit generation order** and module status.
- `docs/ARCH_MAP/<module>.md` — per‑module constraints and implementation record.

If any `docs/ARCH_MAP/*` summary conflicts with canon above, **canon wins** and `docs/ARCH_MAP/*` must be updated to match reality.

---

## 2) Non‑negotiables (must never drift)

- **Authoritative Go server**, thin browser client.
- **Transport:** WSS + binary frames + protobuf payloads (no JSON on hot path).
- **Overworld:** ~10 Hz tick, grid AOI, delta replication, no-change suppression.
- **Battles:** instanced, deterministic lockstep; clients send inputs; server sends **Outcome Timeline** events.
- **Capture-only chess:** pieces are on-board or captured. **No HP, no damage, no healing** anywhere.
- **Determinism:** stable iteration ordering; deterministic RNG where used; RNG outcomes must be represented explicitly in timeline events.
- **Performance:** minimize allocations in hot loops; bounded queues; stable IDs; avoid O(N entities × N clients) scans.

---

## 3) “No silent changes” rule (ledger discipline)

If you must choose something that is not explicitly dictated by canon files:
- Append a new entry to `docs/DECISION_LEDGER.md` (append-only).
- Include:
  - Decision ID (next sequential)
  - What changed
  - Why
  - Paths affected
- Update any relevant `docs/ARCH_MAP/*` files to reflect the choice.

Never “fix” protocol IDs or enum numbers unless the change is explicitly ledgered **and** mirrored in:
- `proto/game.proto`
- `internal/protocol/*`
- any docs that reference the IDs

---

## 4) How to choose work each run (ARCH_MAP pipeline)

Always follow `docs/ARCH_MAP/README.md` **Generation Order**.

Selection rules:
1) Resume the first module marked **IN_PROGRESS**.
2) Else implement the first module marked **TODO**.
3) If all are **DONE**, stop.

For the selected module:
- Open `docs/ARCH_MAP/<module>.md` and obey its constraints.
- The module sheet’s YAML frontmatter defines:
  - `owner:` (allowed primary directory)
  - `touchpoints:` (explicit external files allowed)
  - `depends_on:` (modules that must already be DONE)

If a required file/path is missing, stop and report:
```
MISSING INPUTS:
- <exact/path>
```
Do not guess or invent.

---

## 5) Scope lock (do not “helpfully refactor”)

Default allowed edits per run:
- Files under the module `owner:` path.
- Only the files listed in that module’s `touchpoints:` (if any).
- Documentation updates:
  - `docs/ARCH_MAP/*` (only to record what was implemented)
  - `docs/STATE_HANDOFF.md` (rewrite for this run)
  - `docs/DECISION_LEDGER.md` (append-only)

Any other change is forbidden unless canon explicitly requires it and it is recorded in the ledger.

---

## 6) Editing style (Codex Web mode)

Prefer **editing files in-repo** using patch/diff tools.
Avoid dumping entire file contents unless explicitly requested.

Before making edits, output a short **CHANGE PLAN**:
- files to create/modify (paths)
- any external touchpoints
- validation commands you will run

After edits, output:
- CHANGED FILES (paths)
- VALIDATION RESULTS
- NEXT MODULE (from README order)

---

## 7) Validation (best effort)

After each run:
- `gofmt -w .`
- `go test ./...`

If the environment prevents running commands, state exactly what failed and why, and perform the next best static check (e.g., `go test` for specific packages, or compile-limited checks).

Do not claim tests passed if they did not run.

---

## 8) Documentation updates required every run

### Update the module sheet: `docs/ARCH_MAP/<module>.md`
- Update YAML frontmatter:
  - `status:` (DONE or IN_PROGRESS)
  - `generated_files:` (exact paths written/updated this run)
  - `touchpoints:` (exact external files changed this run)
  - `last_updated:` (YYYY-MM-DD)
- Update body sections:
  - Generated/Modified Files
  - Interfaces / Contracts (public types/functions, msg types used)
  - Algorithmic Invariants Implemented
  - Remaining Work (empty if DONE)

### Update `docs/ARCH_MAP/README.md`
- Flip module status markers accordingly.

### Rewrite `docs/STATE_HANDOFF.md`
Include:
- what changed (paths)
- decisions appended (if any)
- how to validate
- next module to work on

---

## 9) Go server implementation constraints

- Keep packages cohesive and small.
- Minimize allocations in hot loops; reuse buffers where reasonable.
- Use stable ordering (e.g., sort by entity_id where required for deterministic diffs).
- Bounded queues for network writes; backpressure rules must match canon:
  - overworld deltas may be dropped/coalesced
  - battle timelines must never be dropped (disconnect if necessary)

---

## 10) Client constraints (when applicable)

- Thin client: render/UI only.
- Correctness lives on server; client may do interpolation and cosmetic input echo only.
- Battle visuals must follow Outcome Timeline strictly.

---

## 11) Stop conditions (do not continue)

Stop immediately (no code) if:
- Canon files are missing.
- Selected module’s sheet is missing.
- Required dependencies/modules are not DONE and are required by `depends_on`.
- Protocol IDs/enums are ambiguous or inconsistent (must be ledgered and reconciled first).

---

_End of AGENTS.md_
