package ws_gateway

import (
	"context"
	"sync"
	"time"

	"github.com/coder/websocket"

	"example.com/mvp-repo/internal/net/frame"
	"example.com/mvp-repo/internal/protocol"
	"example.com/mvp-repo/internal/router"
)

type conn struct {
	ws         *deadlineConn
	router     *router.Router
	cfg        Config
	queue      *outboundQueue
	notifyCh   chan struct{}
	pool       *sync.Pool
	remoteAddr string
	bound      bool
	playerID   uint64
}

type deadlineConn struct {
	conn          *websocket.Conn
	readDeadline  time.Time
	writeDeadline time.Time
}

func (c *deadlineConn) SetReadDeadline(deadline time.Time) error {
	c.readDeadline = deadline
	return nil
}

func (c *deadlineConn) SetWriteDeadline(deadline time.Time) error {
	c.writeDeadline = deadline
	return nil
}

func (c *deadlineConn) Read(ctx context.Context) (websocket.MessageType, []byte, error) {
	if c.readDeadline.IsZero() {
		return c.conn.Read(ctx)
	}
	readCtx, cancel := context.WithDeadline(ctx, c.readDeadline)
	defer cancel()
	return c.conn.Read(readCtx)
}

func (c *deadlineConn) Write(ctx context.Context, typ websocket.MessageType, data []byte) error {
	if c.writeDeadline.IsZero() {
		return c.conn.Write(ctx, typ, data)
	}
	writeCtx, cancel := context.WithDeadline(ctx, c.writeDeadline)
	defer cancel()
	return c.conn.Write(writeCtx, typ, data)
}

func (c *deadlineConn) Close(code websocket.StatusCode, reason string) error {
	return c.conn.Close(code, reason)
}

func newConn(ws *websocket.Conn, router *router.Router, cfg Config, pool *sync.Pool, remoteAddr string) *conn {
	return &conn{
		ws:         &deadlineConn{conn: ws},
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
		if err := c.ws.SetReadDeadline(time.Now().Add(c.cfg.ReadTimeout)); err != nil {
			return err
		}
		msgType, data, err := c.ws.Read(ctx)
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
		if err := c.ws.SetWriteDeadline(time.Now().Add(c.cfg.WriteTimeout)); err != nil {
			return err
		}
		err := c.ws.Write(ctx, websocket.MessageBinary, frameBytes)
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
