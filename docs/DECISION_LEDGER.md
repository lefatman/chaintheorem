# DECISION LEDGER

Purpose: track **explicit** decisions and any future changes without drift.

## Format rules
- Append-only.
- No silent changes.
- If something must change later, add an **AMENDMENT** entry; do not rewrite history.

### Decision entry template
DECISION ####: <title>
- Date: YYYY-MM-DD
- Status: LOCKED | AMENDED
- Context:
- Options:
- Decision:
- Why:
- Impact:

### Amendment format (hard)
AMENDMENT ####: CHANGED FROM → TO
- Date: YYYY-MM-DD
- Why:
- Impact:
- Follow-ups:

---

## Decisions (seed)

DECISION 0001: WebSocket library for Go server (stdlib-only is impossible for WS)
- Date: 2026-02-02
- Status: LOCKED
- Context: We need WSS and RFC6455 WebSocket support; Go stdlib does not provide a full WebSocket server API.
- Options:
  - github.com/gorilla/websocket
  - github.com/coder/websocket (maintained successor of nhooyr/websocket)
- Decision: Use `github.com/coder/websocket`
- Why:
  - Minimal, idiomatic API; integrates cleanly with net/http upgrade flow.
  - Actively maintained as the continuation of nhooyr/websocket.
- Impact:
  - Server WS gateway package will wrap coder/websocket Conn.
  - Framing + protobuf payload rules remain owned by protocol contract.

DECISION 0002: Protobuf generation approach
- Date: 2026-02-02
- Status: LOCKED
- Context: Protocol contract requires protobuf messages; we need repeatable generation.
- Options:
  - protoc + protoc-gen-go
  - buf (deferred; adds tooling surface)
- Decision: Use `protoc` with `protoc-gen-go` (google.golang.org/protobuf)
- Why:
  - Lowest complexity and broadly supported.
- Impact:
  - Scripts added in `scripts/` and documented in `proto/README.md`.
  - Generated Go output location is fixed (see `proto/README.md`).

DECISION 0003: Client bundler approach
- Date: 2026-02-02
- Status: LOCKED
- Context: PixiJS client needs modern dev/build tooling.
- Options:
  - Vite
  - Webpack
- Decision: Use Vite
- Why:
  - PixiJS ecosystem supports/recommends Vite templates; fast iteration.
- Impact:
  - `client/` will be a Vite project in a later batch (not created yet).

DECISION 0004: AOI defaults (not specified explicitly in canon)
- Date: 2026-02-02
- Status: LOCKED
- Context: AOI doc defines grid cells + radius R but not exact values.
- Decision:
  - `cell_size_tiles = 16`
  - `radius_cells = 1`  (watch 3x3 cells around current cell)
- Why:
  - Minimal, sensible starting scope; easy to reason about; adjustable without protocol changes.
- Impact:
  - Reflected in `config/server.json`.
  - If changed later, must be amended here (not silent).

DECISION 0005: Protobuf schemas for PING/PONG/ERROR (msg types present; schema omitted in sketch)
- Date: 2026-02-02
- Status: LOCKED
- Context: Protocol contract defines msg_type IDs for PING/PONG/ERROR but does not specify protobuf payload shapes.
- Options:
  - Empty messages for keepalive; minimal (Ping{}, Pong{})
  - Error as empty message
  - Error with fields (code + text)
- Decision:
  - PING and PONG are empty protobuf messages.
  - ERROR is `message Error { uint32 code = 1; string text = 2; }`
- Why:
  - Keepalive needs no payload; empty messages are smallest and unambiguous.
  - Error benefits from carrying a stable code plus human-readable text while remaining minimal.
- Impact:
  - Implemented in `proto/game.proto`.
  - Error codes remain implementation-defined until canon formalizes them.

DECISION 0006: PieceType numeric IDs (not specified in protocol contract)
- Date: 2026-02-02
- Status: LOCKED
- Context: Protocol sketch references `PieceType` (e.g., promotion) but does not define canonical numeric IDs.
- Options:
  - Leave as raw uint32 everywhere; avoid an enum
  - Define an internal canonical mapping aligned with common chess ordering
- Decision:
  - Internal mapping (used by code, not locked into protobuf fields yet):
    - 0=UNSPEC, 1=PAWN, 2=KNIGHT, 3=BISHOP, 4=ROOK, 5=QUEEN, 6=KING
- Why:
  - Common, predictable ordering; matches config/gameplay piece ordering; keeps proto flexible (promotion is still uint32).
- Impact:
  - Implemented in `internal/protocol/enums.go`.
  - If canon later specifies IDs, amend this decision and update code accordingly.

DECISION 0007: Dialect-specific migrations for SQLite and Postgres
- Date: 2026-02-02
- Status: LOCKED
- Context: Persist requires schema portability across SQLite (dev) and Postgres (prod) with type differences and placeholder syntax.
- Options:
  - Single shared migration set using only lowest-common-denominator SQL types
  - Separate per-dialect migrations with identical logical schema
- Decision:
  - Maintain separate embedded migration sets for SQLite and Postgres, keeping schema parity while using dialect-correct types.
- Why:
  - SQLite and Postgres differ on identity/bytea types; separate files avoid unsafe compromises while keeping deterministic order.
- Impact:
  - `internal/persist/migrate.go` selects migrations by dialect.
  - `internal/persist/migrations/sqlite/` and `internal/persist/migrations/postgres/` added.

DECISION 0008: Loadout per-piece ability columns
- Date: 2026-02-02
- Status: LOCKED
- Context: Canon schema specifies `ability_*piece` columns without explicit column names.
- Options:
  - Store a serialized list of per-piece assignments
  - Expand to fixed columns per piece type
- Decision:
  - Use fixed columns `ability_pawn`, `ability_knight`, `ability_bishop`, `ability_rook`, `ability_queen`, `ability_king`.
- Why:
  - Matches canonical piece types and keeps storage deterministic and query-friendly.
- Impact:
  - Implemented in `internal/persist/migrations/*/001_init.sql`.
  - Used by `internal/persist/loadouts_repo.go`.

DECISION 0009: Driver registration owned by callers
- Date: 2026-02-02
- Status: LOCKED
- Context: Persist package must support SQLite/Postgres but canon does not mandate driver choice.
- Options:
  - Persist package imports concrete drivers
  - Persist package accepts driver name and expects caller registration
- Decision:
  - Persist package uses `database/sql` and requires callers to register driver imports.
- Why:
  - Keeps persist layer stdlib-first and avoids pinning driver choices prematurely.
- Impact:
  - `internal/persist/db.go` config uses Driver/DSN and does not import drivers.
