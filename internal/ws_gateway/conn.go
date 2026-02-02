package ws_gateway

import (
	"context"
	"sync"

	"github.com/coder/websocket"

	"example.com/mvp-repo/internal/net/frame"
	"example.com/mvp-repo/internal/protocol"
	"example.com/mvp-repo/internal/router"
)

type conn struct {
	ws         *websocket.Conn
	router     *router.Router
	cfg        Config
	queue      *outboundQueue
	notifyCh   chan struct{}
	pool       *sync.Pool
	remoteAddr string
	bound      bool
	playerID   uint64
}

func newConn(ws *websocket.Conn, router *router.Router, cfg Config, pool *sync.Pool, remoteAddr string) *conn {
	return &conn{
		ws:         ws,
		router:     router,
		cfg:        cfg,
		queue:      newOutboundQueue(cfg.WriteQueueMaxFrames),
		notifyCh:   make(chan struct{}, 1),
		pool:       pool,
		remoteAddr: remoteAddr,
	}
}

func (c *conn) run(ctx context.Context) error {
	errCh := make(chan error, 2)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		errCh <- c.readLoop(ctx)
	}()
	go func() {
		errCh <- c.writeLoop(ctx)
	}()

	err := <-errCh
	cancel()
	c.queue.Close()
	_ = c.ws.Close(websocket.StatusNormalClosure, "closing")
	<-errCh
	return err
}

func (c *conn) readLoop(ctx context.Context) error {
	for {
		readCtx, cancel := context.WithTimeout(ctx, c.cfg.ReadTimeout)
		msgType, data, err := c.ws.Read(readCtx)
		cancel()
		if err != nil {
			return err
		}
		if msgType != websocket.MessageBinary {
			return ErrUnsupportedFrame
		}
		wireType, payload, err := frame.Decode(data, c.cfg.ReadLimitBytes)
		if err != nil {
			return err
		}
		msg := protocol.MsgType(wireType)

		if msg == protocol.MSG_PING {
			if err := c.Send(protocol.MSG_PONG, nil); err != nil {
				return err
			}
			continue
		}

		if !c.bound {
			if msg != protocol.MSG_HELLO {
				return ErrUnauthenticated
			}
			ctx := router.Context{
				PlayerID:   c.playerID,
				RemoteAddr: c.remoteAddr,
				Sender:     c,
			}
			if err := c.router.Dispatch(ctx, msg, payload); err != nil {
				return err
			}
			c.bound = true
			continue
		}

		rctx := router.Context{
			PlayerID:   c.playerID,
			RemoteAddr: c.remoteAddr,
			Sender:     c,
		}
		if err := c.router.Dispatch(rctx, msg, payload); err != nil {
			return err
		}
	}
}

func (c *conn) writeLoop(ctx context.Context) error {
	for {
		frameBytes := c.queue.Next()
		if frameBytes == nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c.notifyCh:
				continue
			}
		}
		writeCtx, cancel := context.WithTimeout(ctx, c.cfg.WriteTimeout)
		err := c.ws.Write(writeCtx, websocket.MessageBinary, frameBytes)
		cancel()
		c.putBuffer(frameBytes)
		if err != nil {
			return err
		}
	}
}

func (c *conn) Send(msgType protocol.MsgType, payload []byte) error {
	buffer := c.getBuffer(len(payload))
	frameBytes, err := frame.Encode(buffer[:0], uint16(msgType), payload, c.cfg.ReadLimitBytes)
	if err != nil {
		c.putBuffer(buffer)
		return err
	}

	if isDroppable(msgType) && c.cfg.OverworldDeltaCoalesce {
		replaced, ok := c.queue.SetDroppable(frameBytes)
		if !ok {
			c.putBuffer(frameBytes)
			return ErrBackpressure
		}
		if replaced != nil {
			c.putBuffer(replaced)
		}
		c.notify()
		return nil
	}

	if !c.queue.Enqueue(frameBytes) {
		c.putBuffer(frameBytes)
		if isDroppable(msgType) {
			return nil
		}
		_ = c.Close("backpressure")
		return ErrBackpressure
	}
	c.notify()
	return nil
}

func (c *conn) Close(reason string) error {
	return c.ws.Close(websocket.StatusPolicyViolation, reason)
}

func (c *conn) notify() {
	select {
	case c.notifyCh <- struct{}{}:
	default:
	}
}

func (c *conn) getBuffer(payloadLen int) []byte {
	ptr := c.pool.Get().(*[]byte)
	buf := *ptr
	need := frame.HeaderLen + payloadLen
	if cap(buf) < need {
		buf = make([]byte, 0, need)
	}
	return buf[:0]
}

func (c *conn) putBuffer(buf []byte) {
	if buf == nil {
		return
	}
	buf = buf[:0]
	c.pool.Put(&buf)
}

func isDroppable(msgType protocol.MsgType) bool {
	switch msgType {
	case protocol.MSG_WORLD_DELTA:
		return true
	default:
		return false
	}
}
