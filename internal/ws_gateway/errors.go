package ws_gateway

import "errors"

var (
	ErrRouterRequired   = errors.New("ws_gateway: router is required")
	ErrBackpressure     = errors.New("ws_gateway: backpressure")
	ErrUnsupportedFrame = errors.New("ws_gateway: unsupported frame")
	ErrUnauthenticated  = errors.New("ws_gateway: unauthenticated")
)
