package router

import "example.com/mvp-repo/internal/protocol"

type AuthHandler interface {
	HandleHello(ctx Context, payload []byte) error
}

type WorldHandler interface {
	HandleMoveIntent(ctx Context, payload []byte) error
}

type ChatHandler interface {
	HandleChatSend(ctx Context, payload []byte) error
}

type BattleHandler interface {
	HandleTurnInput(ctx Context, payload []byte) error
}

func (r *Router) RegisterAuth(handler AuthHandler) {
	if handler == nil {
		return
	}
	r.Register(protocol.MSG_HELLO, handler.HandleHello)
}

func (r *Router) RegisterWorld(handler WorldHandler) {
	if handler == nil {
		return
	}
	r.Register(protocol.MSG_WORLD_MOVE_INTENT, handler.HandleMoveIntent)
}

func (r *Router) RegisterChat(handler ChatHandler) {
	if handler == nil {
		return
	}
	r.Register(protocol.MSG_CHAT_SEND, handler.HandleChatSend)
}

func (r *Router) RegisterBattle(handler BattleHandler) {
	if handler == nil {
		return
	}
	r.Register(protocol.MSG_BATTLE_TURN_INPUT, handler.HandleTurnInput)
}
