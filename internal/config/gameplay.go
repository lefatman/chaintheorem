package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type GameplayConfig struct {
	SchemaVersion int               `json:"schema_version"`
	Canon         GameplayCanon     `json:"canon"`
	PieceTypes    []PieceTypeConfig `json:"piece_types"`
	Elements      []ElementConfig   `json:"elements"`
	AbilitySets   AbilitySetsConfig `json:"ability_sets"`
	Abilities     []AbilityConfig   `json:"abilities"`
	Items         []ItemConfig      `json:"items"`
}

type GameplayCanon struct {
	ElementsCount  int      `json:"elements_count"`
	ItemsCount     int      `json:"items_count"`
	AbilitiesCount int      `json:"abilities_count"`
	Notes          []string `json:"notes"`
}

type PieceTypeConfig struct {
	Key  string `json:"key"`
	Rank int    `json:"rank"`
}

type ElementConfig struct {
	ID       int            `json:"id"`
	Key      string         `json:"key"`
	Name     string         `json:"name"`
	Passives map[string]any `json:"passives"`
}

type AbilitySetsConfig struct {
	DefensiveAbilityIDs       []int `json:"defensive_ability_ids"`
	OffensiveAbilityIDs       []int `json:"offensive_ability_ids"`
	RemoteOffensiveCaptureIDs []int `json:"remote_offensive_capture_ability_ids"`
}

type AbilityConfig struct {
	ID         int            `json:"id"`
	Key        string         `json:"key"`
	Name       string         `json:"name"`
	Scope      string         `json:"scope"`
	Category   string         `json:"category"`
	Consumable bool           `json:"consumable"`
	Charges    map[string]any `json:"charges"`
	Rules      map[string]any `json:"rules"`
}

type ItemConfig struct {
	ID                  int            `json:"id"`
	Key                 string         `json:"key"`
	Name                string         `json:"name"`
	SlotCost            int            `json:"slot_cost"`
	Effects             map[string]any `json:"effects"`
	IncompatibleItemIDs []int          `json:"incompatible_item_ids"`
}

func LoadGameplayConfig(path string) (GameplayConfig, error) {
	var cfg GameplayConfig
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
		return cfg, errors.New("gameplay config: unexpected trailing data")
	}
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (cfg GameplayConfig) Validate() error {
	if cfg.SchemaVersion != 1 {
		return fmt.Errorf("gameplay config: schema_version must be 1")
	}
	if cfg.Canon.ElementsCount <= 0 || cfg.Canon.ItemsCount <= 0 || cfg.Canon.AbilitiesCount <= 0 {
		return fmt.Errorf("gameplay config: canon counts must be > 0")
	}
	if len(cfg.Elements) != cfg.Canon.ElementsCount {
		return fmt.Errorf("gameplay config: elements count mismatch")
	}
	if len(cfg.Items) != cfg.Canon.ItemsCount {
		return fmt.Errorf("gameplay config: items count mismatch")
	}
	if len(cfg.Abilities) != cfg.Canon.AbilitiesCount {
		return fmt.Errorf("gameplay config: abilities count mismatch")
	}
	if len(cfg.PieceTypes) == 0 {
		return fmt.Errorf("gameplay config: piece_types must be non-empty")
	}
	if err := validateSequentialIDs(cfg.Elements, cfg.Canon.ElementsCount, 0, "elements"); err != nil {
		return err
	}
	if err := validateSequentialIDs(cfg.Items, cfg.Canon.ItemsCount, 1, "items"); err != nil {
		return err
	}
	if err := validateSequentialIDs(cfg.Abilities, cfg.Canon.AbilitiesCount, 1, "abilities"); err != nil {
		return err
	}
	return nil
}

type idProvider interface {
	getID() int
}

func validateSequentialIDs[T idProvider](items []T, count int, minID int, label string) error {
	seen := make([]bool, count)
	maxID := minID + count - 1
	for _, item := range items {
		id := item.getID()
		if id < minID || id > maxID {
			return fmt.Errorf("gameplay config: %s id %d out of range", label, id)
		}
		index := id - minID
		if seen[index] {
			return fmt.Errorf("gameplay config: %s id %d duplicated", label, id)
		}
		seen[index] = true
	}
	for i, ok := range seen {
		if !ok {
			return fmt.Errorf("gameplay config: %s id %d missing", label, i+minID)
		}
	}
	return nil
}

func (e ElementConfig) getID() int { return e.ID }
func (a AbilityConfig) getID() int { return a.ID }
func (i ItemConfig) getID() int    { return i.ID }
