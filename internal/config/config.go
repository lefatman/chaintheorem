package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type ServerConfig struct {
	SchemaVersion int               `json:"schema_version"`
	HTTP          HTTPConfig        `json:"http"`
	WS            WSConfig          `json:"ws"`
	Overworld     OverworldConfig   `json:"overworld"`
	Battle        BattleConfig      `json:"battle"`
	Auth          AuthConfig        `json:"auth"`
	Persistence   PersistenceConfig `json:"persistence"`
}

type HTTPConfig struct {
	ListenAddr string `json:"listen_addr"`
}

type WSConfig struct {
	ListenAddr             string    `json:"listen_addr"`
	Path                   string    `json:"path"`
	ReadLimitBytes         uint32    `json:"read_limit_bytes"`
	WriteQueueMaxFrames    int       `json:"write_queue_max_frames"`
	OverworldDeltaCoalesce bool      `json:"overworld_delta_coalesce"`
	TLS                    TLSConfig `json:"tls"`
}

type TLSConfig struct {
	Enabled  bool   `json:"enabled"`
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

type OverworldConfig struct {
	TickHz      int               `json:"tick_hz"`
	GridAOI     GridAOIConfig     `json:"grid_aoi"`
	Replication ReplicationConfig `json:"replication"`
}

type GridAOIConfig struct {
	CellSizeTiles int `json:"cell_size_tiles"`
	RadiusCells   int `json:"radius_cells"`
}

type ReplicationConfig struct {
	NoChangeSuppression                bool   `json:"no_change_suppression"`
	StableSort                         string `json:"stable_sort"`
	MaxPendingOverworldDeltasPerClient int    `json:"max_pending_overworld_deltas_per_client"`
	OnBackpressure                     string `json:"on_backpressure"`
}

type BattleConfig struct {
	DeterministicLockstep               bool      `json:"deterministic_lockstep"`
	OutcomeTimelineIsOnlyAnimationTruth bool      `json:"outcome_timeline_is_only_animation_truth"`
	HistoryPliesForRedo                 int       `json:"history_plies_for_redo"`
	RNG                                 RNGConfig `json:"rng"`
}

type RNGConfig struct {
	PRNG string `json:"prng"`
}

type AuthConfig struct {
	SessionTokenBytes int `json:"session_token_bytes"`
	SessionTTLSeconds int `json:"session_ttl_seconds"`
}

type PersistenceConfig struct {
	DevDriver  string `json:"dev_driver"`
	ProdDriver string `json:"prod_driver"`
}

func LoadServerConfig(path string) (ServerConfig, error) {
	var cfg ServerConfig
	file, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&cfg); err != nil {
		return cfg, err
	}
	if dec.More() {
		return cfg, errors.New("server config: unexpected trailing data")
	}
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (cfg ServerConfig) Validate() error {
	if cfg.SchemaVersion != 1 {
		return fmt.Errorf("server config: schema_version must be 1")
	}
	if cfg.HTTP.ListenAddr == "" {
		return fmt.Errorf("server config: http.listen_addr is required")
	}
	if cfg.WS.ListenAddr == "" {
		return fmt.Errorf("server config: ws.listen_addr is required")
	}
	if cfg.WS.Path == "" || cfg.WS.Path[0] != '/' {
		return fmt.Errorf("server config: ws.path must start with '/'")
	}
	if cfg.WS.ReadLimitBytes == 0 {
		return fmt.Errorf("server config: ws.read_limit_bytes must be > 0")
	}
	if cfg.WS.WriteQueueMaxFrames <= 0 {
		return fmt.Errorf("server config: ws.write_queue_max_frames must be > 0")
	}
	if cfg.WS.TLS.Enabled {
		if cfg.WS.TLS.CertFile == "" || cfg.WS.TLS.KeyFile == "" {
			return fmt.Errorf("server config: ws.tls.cert_file and ws.tls.key_file are required when tls.enabled")
		}
	}
	if cfg.Overworld.TickHz <= 0 {
		return fmt.Errorf("server config: overworld.tick_hz must be > 0")
	}
	if cfg.Overworld.GridAOI.CellSizeTiles <= 0 {
		return fmt.Errorf("server config: overworld.grid_aoi.cell_size_tiles must be > 0")
	}
	if cfg.Overworld.GridAOI.RadiusCells <= 0 {
		return fmt.Errorf("server config: overworld.grid_aoi.radius_cells must be > 0")
	}
	if cfg.Overworld.Replication.MaxPendingOverworldDeltasPerClient <= 0 {
		return fmt.Errorf("server config: overworld.replication.max_pending_overworld_deltas_per_client must be > 0")
	}
	if cfg.Auth.SessionTokenBytes <= 0 {
		return fmt.Errorf("server config: auth.session_token_bytes must be > 0")
	}
	if cfg.Auth.SessionTTLSeconds <= 0 {
		return fmt.Errorf("server config: auth.session_ttl_seconds must be > 0")
	}
	if cfg.Persistence.DevDriver == "" || cfg.Persistence.ProdDriver == "" {
		return fmt.Errorf("server config: persistence.dev_driver and persistence.prod_driver are required")
	}
	return nil
}
