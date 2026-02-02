package ws_gateway

import (
	"fmt"
	"time"
)

const (
	DefaultReadTimeout  = 30 * time.Second
	DefaultWriteTimeout = 10 * time.Second
)

type Config struct {
	ReadLimitBytes         uint32
	WriteQueueMaxFrames    int
	OverworldDeltaCoalesce bool
	ReadTimeout            time.Duration
	WriteTimeout           time.Duration
}

func (cfg Config) Validate() error {
	if cfg.ReadLimitBytes == 0 {
		return fmt.Errorf("ws_gateway: ReadLimitBytes must be > 0")
	}
	if cfg.WriteQueueMaxFrames <= 0 {
		return fmt.Errorf("ws_gateway: WriteQueueMaxFrames must be > 0")
	}
	if cfg.ReadTimeout <= 0 {
		return fmt.Errorf("ws_gateway: ReadTimeout must be > 0")
	}
	if cfg.WriteTimeout <= 0 {
		return fmt.Errorf("ws_gateway: WriteTimeout must be > 0")
	}
	return nil
}
