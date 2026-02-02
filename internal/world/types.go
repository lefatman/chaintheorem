// File: internal/world/types.go
package world

import "errors"

var (
	ErrInvalidIntent = errors.New("world: invalid move intent")
	ErrUnknownEntity = errors.New("world: unknown entity")
)

type Entity struct {
	ID   uint64
	X    int32
	Y    int32
	Kind uint32
}

type MoveIntent struct {
	DX int32
	DY int32
}

func (m MoveIntent) Valid() bool {
	if m.DX < -1 || m.DX > 1 || m.DY < -1 || m.DY > 1 {
		return false
	}
	if m.DX != 0 && m.DY != 0 {
		return false
	}
	return true
}

type EntitySource interface {
	AppendEntities(dst []Entity) []Entity
	EntityByID(id uint64) (Entity, bool)
}

type MoveIntentSink interface {
	SetMoveIntent(playerID uint64, dx, dy int32) error
}
