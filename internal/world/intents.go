// File: internal/world/intents.go
package world

import "sort"

const invalidEntityIndex = -1

type IntentStore struct {
	playerIDs   []uint64
	entityIDs   []uint64
	entityIndex []int
	intents     []MoveIntent
	active      []bool
	index       map[uint64]int
}

func NewIntentStore(capacity int) *IntentStore {
	if capacity < 0 {
		capacity = 0
	}
	return &IntentStore{
		playerIDs:   make([]uint64, 0, capacity),
		entityIDs:   make([]uint64, 0, capacity),
		entityIndex: make([]int, 0, capacity),
		intents:     make([]MoveIntent, 0, capacity),
		active:      make([]bool, 0, capacity),
		index:       make(map[uint64]int, capacity),
	}
}

func (s *IntentStore) SetMoveIntent(playerID uint64, dx, dy int32) error {
	intent := MoveIntent{DX: dx, DY: dy}
	if !intent.Valid() {
		return ErrInvalidIntent
	}
	idx := s.ensurePlayer(playerID)
	s.intents[idx] = intent
	s.active[idx] = true
	return nil
}

func (s *IntentStore) BindPlayerEntity(playerID uint64, entityID uint64, store *Store) error {
	if store == nil {
		return ErrUnknownEntity
	}
	idx, ok := store.EntityIndex(entityID)
	if !ok {
		return ErrUnknownEntity
	}
	playerIdx := s.ensurePlayer(playerID)
	s.entityIDs[playerIdx] = entityID
	s.entityIndex[playerIdx] = idx
	return nil
}

func (s *IntentStore) ensurePlayer(playerID uint64) int {
	if idx, ok := s.index[playerID]; ok {
		return idx
	}
	insertAt := sort.Search(len(s.playerIDs), func(i int) bool {
		return s.playerIDs[i] >= playerID
	})
	s.playerIDs = append(s.playerIDs, 0)
	s.entityIDs = append(s.entityIDs, 0)
	s.entityIndex = append(s.entityIndex, 0)
	s.intents = append(s.intents, MoveIntent{})
	s.active = append(s.active, false)

	if insertAt < len(s.playerIDs)-1 {
		copy(s.playerIDs[insertAt+1:], s.playerIDs[insertAt:])
		copy(s.entityIDs[insertAt+1:], s.entityIDs[insertAt:])
		copy(s.entityIndex[insertAt+1:], s.entityIndex[insertAt:])
		copy(s.intents[insertAt+1:], s.intents[insertAt:])
		copy(s.active[insertAt+1:], s.active[insertAt:])
		for i := insertAt + 1; i < len(s.playerIDs); i++ {
			s.index[s.playerIDs[i]] = i
		}
	}

	s.playerIDs[insertAt] = playerID
	s.entityIDs[insertAt] = 0
	s.entityIndex[insertAt] = invalidEntityIndex
	s.intents[insertAt] = MoveIntent{}
	s.active[insertAt] = false
	s.index[playerID] = insertAt
	return insertAt
}
