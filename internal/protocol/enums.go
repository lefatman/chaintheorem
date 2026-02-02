package protocol

// ElementId values are canonical (match Protocol Contract tables / proto enums).
type ElementId uint8

const (
	ELEMENT_WATER     ElementId = 0
	ELEMENT_FIRE      ElementId = 1
	ELEMENT_EARTH     ElementId = 2
	ELEMENT_AIR_WIND  ElementId = 3
	ELEMENT_LIGHTNING ElementId = 4
)

// AbilityId values are canonical (match Protocol Contract tables / proto enums).
type AbilityId uint16

const (
	ABILITY_BLOCK_PATH   AbilityId = 1
	ABILITY_STALWART      AbilityId = 2
	ABILITY_BELLIGERENT   AbilityId = 3
	ABILITY_REDO          AbilityId = 4
	ABILITY_DOUBLE_KILL   AbilityId = 5
	ABILITY_QUANTUM_KILL  AbilityId = 6
	ABILITY_CHAIN_KILL    AbilityId = 7
	ABILITY_NECROMANCER   AbilityId = 8
)

// ItemId values are canonical (match Protocol Contract tables / proto enums).
type ItemId uint16

const (
	ITEM_MULTITASKERS_SCHEDULE ItemId = 1
	ITEM_POISONED_DAGGER       ItemId = 2
	ITEM_DUAL_ADEPTS_GLOVES    ItemId = 3
	ITEM_TRIPLE_ADEPTS_GLOVES  ItemId = 4
	ITEM_HEADMASTER_RING       ItemId = 5
	ITEM_POT_OF_HUNGER         ItemId = 6
	ITEM_SOLAR_NECKLACE        ItemId = 7
)

// PieceType numeric IDs were not explicitly specified in the Protocol Contract.
// Mapping is ledgered in DECISION 0006 and may be amended if canon later defines IDs.
type PieceType uint8

const (
	PIECE_UNSPEC PieceType = 0
	PIECE_PAWN   PieceType = 1
	PIECE_KNIGHT PieceType = 2
	PIECE_BISHOP PieceType = 3
	PIECE_ROOK   PieceType = 4
	PIECE_QUEEN  PieceType = 5
	PIECE_KING   PieceType = 6
)

// Rank values are canonical chess ranks used by gameplay rules.
type Rank uint8

const (
	RANK_PAWN   Rank = 1
	RANK_KNIGHT Rank = 2
	RANK_BISHOP Rank = 2
	RANK_ROOK   Rank = 3
	RANK_QUEEN  Rank = 4
	RANK_KING   Rank = 5
)

func RankOfPieceType(pt PieceType) Rank {
	switch pt {
	case PIECE_PAWN:
		return RANK_PAWN
	case PIECE_KNIGHT:
		return RANK_KNIGHT
	case PIECE_BISHOP:
		return RANK_BISHOP
	case PIECE_ROOK:
		return RANK_ROOK
	case PIECE_QUEEN:
		return RANK_QUEEN
	case PIECE_KING:
		return RANK_KING
	default:
		return 0
	}
}

type Dir4 uint8

const (
	DIR_N Dir4 = 0
	DIR_E Dir4 = 1
	DIR_S Dir4 = 2
	DIR_W Dir4 = 3
)

type BattleActionType uint8

const (
	BAT_ACT_MOVE      BattleActionType = 0
	BAT_ACT_CHAIN_KILL BattleActionType = 1
)

type TimelineEventType uint8

const (
	EV_MOVE            TimelineEventType = 0
	EV_CAPTURE         TimelineEventType = 1
	EV_EXTRA_CAPTURE   TimelineEventType = 2
	EV_BLOCK_PATH_SET  TimelineEventType = 3
	EV_ABILITY_FIZZLE  TimelineEventType = 4
	EV_REDO_REWIND     TimelineEventType = 5
	EV_PIECE_RESTORED  TimelineEventType = 6
	EV_MATCH_STATE     TimelineEventType = 7
)
