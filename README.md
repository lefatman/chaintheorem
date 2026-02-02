# MVP Repo (Authoritative Go Server + Thin PixiJS Client)

This repo is being generated in **canon-locked batches**.

## Canon (do not violate)
- Go server is authoritative; client is a thin renderer/UI.
- HTTPS for auth/assets; **WSS** for realtime.
- Binary frames + Protobuf payloads; **no JSON on the hot path**.
- Overworld is **10 Hz** tick with **grid AOI** + **delta diffs** + **no-change suppression**.
- Battles are instanced **deterministic lockstep**: inputs in, **Outcome Timeline** out.
- **Capture-only** semantics: pieces are on-board or captured. **No HP / damage / heal** anywhere.

See: `docs/CANON_LOCK.md`.

## What exists in Batch 00
- Repo skeleton + canon lock + decision ledger
- `config/gameplay.json` (elements/abilities/items/loadout rules; canon IDs)
- `config/server.json` (10 Hz + AOI defaults)
- Protobuf generation scripts (no .proto files yet)

## How to run (placeholders — implemented in later batches)
### Server (later)
```bash
# go build ./cmd/server
# ./server -config ./config/server.json
```

### Client (later)
```bash
# cd client
# npm install
# npm run dev
```

## Protobuf generation
- Put `.proto` files under `proto/` (later batch will add them).
- Generate Go types into `internal/proto/gen`:
```bash
make proto
```

## Rule of the land
If anything must change, it must be logged as an explicit amendment in `docs/DECISION_LEDGER.md`.
