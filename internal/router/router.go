package router

import (
	"errors"

	"example.com/mvp-repo/internal/protocol"
)

const maxMsgType = 1 << 16

var (
	ErrUnhandled       = errors.New("router: unhandled msg_type")
	ErrUnauthenticated = errors.New("router: unauthenticated")
	ErrUnimplemented   = errors.New("router: unimplemented")
)

type Sender interface {
	Send(msgType protocol.MsgType, payload []byte) error
	Close(reason string) error
}

type Context struct {
	PlayerID   uint64
	RemoteAddr string
	Sender     Sender
}

type Handler func(ctx Context, payload []byte) error

type Router struct {
	handlers [maxMsgType]Handler
}

func New() *Router {
	return &Router{}
}

func (r *Router) Register(msgType protocol.MsgType, handler Handler) {
	r.handlers[msgType] = handler
}

func (r *Router) Dispatch(ctx Context, msgType protocol.MsgType, payload []byte) error {
	if int(msgType) >= len(r.handlers) {
		return ErrUnhandled
	}
	h := r.handlers[msgType]
	if h == nil {
		return ErrUnhandled
	}
	return h(ctx, payload)
}
