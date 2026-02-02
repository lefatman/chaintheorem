# Global Contracts and Invariants

This document is the “gravity” that prevents drift.

## Canon binding
- `proto/game.proto` and `internal/protocol/*` are canonical for wire IDs and enums.
- msg_type numeric IDs and enum IDs are binding; never renumber.

## Wire + simulation invariants
- **Realtime is WSS binary frames** carrying **Protobuf payloads**.
- Frame header is **little-endian**:
  - `u16 msg_type`
  - `u32 payload_len`
  - `payload_len` bytes payload
- **No JSON on hot path** (realtime). JSON is acceptable only on HTTPS endpoints.
- **Overworld**: authoritative **10 Hz** tick; AOI is **grid-based**; deltas are **suppressed if no change**.
- **Battles**: instanced **deterministic lockstep**:
  - clients send **inputs**
  - server returns **Outcome Timeline** events (only animation truth)
- **Capture-only semantics**: no HP/damage/heal anywhere.

## Decision bindings
- **DECISION 0005**: PING and PONG are empty protobuf messages; ERROR is `message Error { uint32 code = 1; string text = 2; }`.
- **DECISION 0006**: PieceType mapping (internal): 0=UNSPEC, 1=PAWN, 2=KNIGHT, 3=BISHOP, 4=ROOK, 5=QUEEN, 6=KING.

## Determinism rules
- Any randomness is deterministic PRNG seeded at battle start.
- Any RNG outcome is emitted as an explicit timeline event.
- Deterministic ordering: do not iterate over Go maps for anything that affects output.

## “No silent changes” protocol
If you must decide something that is not specified by the canon docs, you must:
1) choose the smallest / simplest option,
2) append a new Decision Ledger entry describing the choice and rationale,
3) ensure downstream modules consume the same decision.

## msg_type registry (canonical)
| msg_type | Direction | Name | Purpose |
| --- | --- | --- | --- |
| 1 | C->S | HELLO | Authenticate WSS connection using session token. |
| 2 | S->C | WELCOME | Bind connection to player_id; initial routing info. |
| 3 | C->S | PING | Keepalive. |
| 4 | S->C | PONG | Keepalive response. |
| 10 | C->S | WORLD_MOVE_INTENT | Request overworld movement (authoritative on server tick). |
| 11 | S->C | WORLD_SNAPSHOT | Full AOI state (initial or resync). |
| 12 | S->C | WORLD_DELTA | AOI delta diff at 10 Hz (no-change suppressed). |
| 20 | C->S | CHAT_SEND | Send chat message. |
| 21 | S->C | CHAT_EVENT | Chat broadcast. |
| 30 | S->C | BATTLE_START | Enter battle instance; includes seed and initial board. |
| 31 | C->S | BATTLE_TURN_INPUT | Submit one ply input for deterministic lockstep. |
| 32 | S->C | BATTLE_OUTCOME_TIMELINE | Authoritative ordered events for one ply. |
| 33 | S->C | BATTLE_END | Terminal battle result + rewards summary. |
| 250 | S->C | ERROR | Error / rejection. |


## Elements, Abilities, Items (canonical tables)
### Elements
| ElementId | Name | Passive rules (canonical) |
| --- | --- | --- |
| 0 | Water | Consumable ability counters are doubled. This doubling is negated when fighting a Lightning army. |
| 1 | Fire | Offensive abilities resolve first. Fire offensive abilities are ineffective against Water armies. |
| 2 | Earth | Remote offensive capture abilities are nullified. This nullification is negated by Fire armies. |
| 3 | Air/Wind | Negates defensive abilities and can move over pieces. Air/Wind passives are negated by Earth armies. |
| 4 | Lightning | Army abilities are slottable at the piece-type level (no Multitasker’s Schedule needed). Abilities have a 50% chance to misfire against Air/Wind armies. |


### Items
| ItemId | Item | Slot cost | Effect |
| --- | --- | --- | --- |
| 1 | Multitasker’s Schedule | 1 | Allows piece-type-level slotting of army abilities (non-Lightning armies). |
| 2 | Poisoned Dagger | 1 | When a piece is captured, removes the capturing piece if the capturer is of lower or equal rank. |
| 3 | Dual Adept’s Gloves | 1 | Adds +1 army-ability slot. |
| 4 | Triple Adept’s Gloves | 2 | Adds +2 army-ability slots; blocks Dual Adept’s Gloves and Headmaster Ring. |
| 5 | Headmaster Ring | 3 | Adds +3 army-ability slots; blocks Dual and Triple Adept’s Gloves. |
| 6 | Pot of Hunger | 1 | Doubles XP gained from winning a match. |
| 7 | Solar Necklace | 1 | Can top up a consumable ability up to 3 times per match. |


### Abilities
Abilities are army-type only. Piece-type slotting is assignment/scoping, not a separate ability class.
| AbilityId | Ability | Scope | Category | Consumable | Canonical effect |
| --- | --- | --- | --- | --- | --- |
| 1 | Block Path | army-wide | defensive | no | After moving, choose a cardinal direction; this piece cannot be captured from that direction until it moves again. |
| 2 | Stalwart | army-wide | defensive | no | Pieces of this type cannot be captured by a lower-rank capturer. |
| 3 | Belligerent | army-wide | defensive | no | Pieces of this type cannot be captured by a higher-rank capturer. |
| 4 | Redo | army-wide | defensive | yes (1 per piece; doubled for Water vs non-Lightning) | When a piece of this type is captured, rewind exactly 2 plies (the capturing ply and the defender ply immediately before it), restoring the full match state to the start of the defender's prior turn. The defender then replays that prior turn with a different move. All effects of those two plies are reversed. Spend 1 Redo charge on the captured piece; the spent charge remains spent after the rewind. |
| 5 | Double Kill | army-wide | offensive | no | On capture, remove one neighboring enemy piece of equal or lower rank; if none exists, no effect. |
| 6 | Quantum Kill | army-wide | offensive | no | On capture, remove one random enemy piece of equal or lower rank. |
| 7 | Chain Kill | army-wide | offensive | no | Active capture: a piece can piggyback on an adjacent allied piece to capture a target as if it were on the ally’s square (remote capture). |
| 8 | Necromancer | army-wide | offensive | yes (side pool; doubled for Water vs non-Lightning) | If your piece captures a higher-rank piece, and you have any captured pieces, return one eligible captured friendly piece to the board at its capture square. |

