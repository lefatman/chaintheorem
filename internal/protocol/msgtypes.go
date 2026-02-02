package protocol

// MsgType values are the canonical u16 values used in the wire frame header.
// These MUST match the Protocol Contract — CONSOLIDATED msg_type table.
type MsgType uint16

const (
	MSG_HELLO MsgType = 1
	MSG_WELCOME MsgType = 2

	MSG_PING MsgType = 3
	MSG_PONG MsgType = 4

	MSG_WORLD_MOVE_INTENT MsgType = 10
	MSG_WORLD_SNAPSHOT    MsgType = 11
	MSG_WORLD_DELTA       MsgType = 12

	MSG_CHAT_SEND  MsgType = 20
	MSG_CHAT_EVENT MsgType = 21

	MSG_BATTLE_START            MsgType = 30
	MSG_BATTLE_TURN_INPUT       MsgType = 31
	MSG_BATTLE_OUTCOME_TIMELINE MsgType = 32
	MSG_BATTLE_END              MsgType = 33

	MSG_ERROR MsgType = 250
)
