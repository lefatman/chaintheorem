package ws_gateway

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/coder/websocket"

	"example.com/mvp-repo/internal/router"
)

type Server struct {
	cfg    Config
	router *router.Router
	pool   sync.Pool
}

func New(cfg Config, router *router.Router) (*Server, error) {
	if router == nil {
		return nil, ErrRouterRequired
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = DefaultReadTimeout
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = DefaultWriteTimeout
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Server{
		cfg:    cfg,
		router: router,
		pool: sync.Pool{
			New: func() any {
				buf := make([]byte, 0, int(cfg.ReadLimitBytes)+8)
				return &buf
			},
		},
	}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	conn.SetReadLimit(int64(s.cfg.ReadLimitBytes))

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	c := newConn(conn, s.router, s.cfg, &s.pool, r.RemoteAddr)
	if err := c.run(ctx); err != nil {
		log.Printf("ws_gateway: disconnect %s: %v", r.RemoteAddr, err)
	}
}
