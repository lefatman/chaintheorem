package app

import (
	"fmt"

	"example.com/mvp-repo/internal/config"
	"example.com/mvp-repo/internal/router"
	"example.com/mvp-repo/internal/ws_gateway"
)

type App struct {
	Router  *router.Router
	Gateway *ws_gateway.Server
}

func New(serverCfg config.ServerConfig) (*App, error) {
	r := router.New()
	auth := authHandler{}
	world := worldHandler{}
	chat := chatHandler{}
	battle := battleHandler{}

	r.RegisterAuth(auth)
	r.RegisterWorld(world)
	r.RegisterChat(chat)
	r.RegisterBattle(battle)

	gwCfg := ws_gateway.Config{
		ReadLimitBytes:         serverCfg.WS.ReadLimitBytes,
		WriteQueueMaxFrames:    serverCfg.WS.WriteQueueMaxFrames,
		OverworldDeltaCoalesce: serverCfg.WS.OverworldDeltaCoalesce,
	}
	gateway, err := ws_gateway.New(gwCfg, r)
	if err != nil {
		return nil, err
	}
	return &App{
		Router:  r,
		Gateway: gateway,
	}, nil
}

type authHandler struct{}

type worldHandler struct{}

type chatHandler struct{}

type battleHandler struct{}

func (authHandler) HandleHello(ctx router.Context, payload []byte) error {
	_ = payload
	if ctx.Sender == nil {
		return fmt.Errorf("app: sender required")
	}
	return nil
}

func (worldHandler) HandleMoveIntent(ctx router.Context, payload []byte) error {
	_ = payload
	if ctx.Sender == nil {
		return fmt.Errorf("app: sender required")
	}
	_ = ctx.Sender.Close("world handler not implemented")
	return router.ErrUnimplemented
}

func (chatHandler) HandleChatSend(ctx router.Context, payload []byte) error {
	_ = payload
	if ctx.Sender == nil {
		return fmt.Errorf("app: sender required")
	}
	_ = ctx.Sender.Close("chat handler not implemented")
	return router.ErrUnimplemented
}

func (battleHandler) HandleTurnInput(ctx router.Context, payload []byte) error {
	_ = payload
	if ctx.Sender == nil {
		return fmt.Errorf("app: sender required")
	}
	_ = ctx.Sender.Close("battle handler not implemented")
	return router.ErrUnimplemented
}
