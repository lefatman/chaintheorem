---
status: todo
owner: internal/loadout
generated_files: []
touchpoints:
  - internal/httpapi
  - internal/persist
  - internal/battle_mgr
  - config/gameplay.json
depends_on:
  - internal_protocol
  - internal_config
  - internal_persist
last_updated: 2026-02-02
---

# internal/loadout — Army Ability Loadout Validation + Assignment

## Purpose
Validate and normalize a player's loadout where:

- **All abilities are army abilities** (single ability class).
- Abilities may be **assigned** either:
  - to **army slots** (default), or
  - to **piece-type slots** (assignment-by-piece-type),
    which is allowed **only if**:
    - the army element is **Lightning**, OR
    - the army has **Multitasker’s Schedule** equipped.

This module prevents illegal assignments, illegal item combos, illegal slot counts, and non-canonical storage forms.

## Ownership
This module owns:
- Loadout validation (slot counts, exclusivity, assignment permissions)
- Loadout normalization (canonical ordering, stable storage format)
- Derived slot computation from element + items

This module must not own:
- DB schema or SQL (owned by `internal/persist`)
- Battle resolution semantics (owned by `internal/battle_engine`)
- Protocol IDs/enums (owned by `proto/` + `internal/protocol`)

## Canon constraints (binding)

### Abilities
- **Single class:** all abilities are “army abilities”.
- There is no separate concept of “piece-type abilities”.
- “Piece-type slotting” is purely **assignment** (scoping) of an army ability to a piece type.

### Slot model (assignment capacity)
- Base slots:
  - **Army assignment slots:** 1
  - **Piece-type assignment slots:** 6 (i.e., up to 6 assignments to specific PieceTypes)
- Maximum army slots = **4**, achieved via items:
  - Dual Adept (+1)
  - Triple Adept (+2)
  - Headmaster’s Tactics (+3)
- Exclusivity (must enforce):
  - Triple blocks Dual and Headmaster
  - Headmaster blocks Dual and Triple

### Permission rule (the critical one)
An army ability may be assigned to a piece-type slot **only if**:
- Element is **Lightning**, OR
- Item **Multitasker’s Schedule** is equipped.

If not permitted, all abilities must be assigned only to army slots.

### IDs
AbilityId / ItemId / ElementId / PieceType numeric IDs must match:
- `proto/game.proto` enums
- `internal/protocol/enums.go`
- `config/gameplay.json`

Never invent or renumber IDs here.

## Loadout data shape (normalized)

Recommended normalized representation:

- element_id
- items: sorted ascending by ItemId (no duplicates)
- ability_assignments: sorted stable list of:
  - AbilityId
  - AssignmentScope:
    - ARMY
    - PIECE_TYPE (includes PieceType)

Constraints:
- An ability may appear at most once per scope target unless canon explicitly allows duplicates (default: no duplicates).
- Total ARMY assignments <= computed army slots.
- Total PIECE_TYPE assignments <= 6.
- PIECE_TYPE assignments forbidden unless permission rule satisfied.

## Algorithms

### ComputeArmySlots(items)
- base = 1
- +1 if Dual Adept
- +2 if Triple Adept
- +3 if Headmaster’s Tactics
- enforce mutual exclusivity among {Dual, Triple, Headmaster}

### CanUsePieceTypeAssignments(element, items)
- return true if element == Lightning OR items contains Multitasker’s Schedule

### ValidateAssignments(loadout)
- Validate item exclusivity
- Compute army slots
- If any PIECE_TYPE assignments exist:
  - require CanUsePieceTypeAssignments == true
- Enforce slot capacities:
  - count(ARMY assignments) <= army_slots
  - count(PIECE_TYPE assignments) <= 6
- Validate IDs exist in config/gameplay.json and match protocol enums
- Normalize ordering for storage

## Integration boundaries (touchpoints)

### internal/httpapi
Optional endpoints (if MVP includes them):
- GET /api/loadout
- POST /api/loadout
POST must re-validate and store the normalized form.

### internal/persist
Stores the normalized representation.
This module does not write SQL.

### internal/battle_mgr
On battle start, validate the selected loadout:
- reject if invalid
- pass normalized assignments forward

## Expected code files (under owner)
- internal/loadout/types.go
  - input loadout + NormalizedLoadout + AbilityAssignment + Scope
- internal/loadout/slots.go
  - ComputeArmySlots + exclusivity checks
- internal/loadout/validate.go
  - ValidateAndNormalize(cfg, input) -> normalized or error
- internal/loadout/normalize.go
  - stable sorting/canonicalization helpers
- internal/loadout/errors.go
  - typed errors suitable for mapping to API responses

## Remaining work
- [ ] Implement internal/loadout package per file list
- [ ] Wire validation into internal/battle_mgr battle start
- [ ] (Optional) Add HTTP endpoints for loadout editing
- [ ] Update status/generated_files when implemented
