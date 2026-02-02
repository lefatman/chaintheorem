// File: internal/world/store.go
package world

type Store struct {
	nextID uint64

	ids   []uint64
	x     []int32
	y     []int32
	kind  []uint32
	alive []bool

	aliveCount int
	index      map[uint64]int
}

func NewStore(capacity int) *Store {
	if capacity < 0 {
		capacity = 0
	}
	return &Store{
		ids:   make([]uint64, 0, capacity),
		x:     make([]int32, 0, capacity),
		y:     make([]int32, 0, capacity),
		kind:  make([]uint32, 0, capacity),
		alive: make([]bool, 0, capacity),
		index: make(map[uint64]int, capacity),
	}
}

func (s *Store) AliveCount() int {
	return s.aliveCount
}

func (s *Store) CreateEntity(x, y int32, kind uint32) uint64 {
	s.nextID++
	id := s.nextID
	idx := len(s.ids)

	s.ids = append(s.ids, id)
	s.x = append(s.x, x)
	s.y = append(s.y, y)
	s.kind = append(s.kind, kind)
	s.alive = append(s.alive, true)
	s.index[id] = idx
	s.aliveCount++
	return id
}

func (s *Store) RemoveEntity(id uint64) bool {
	idx, ok := s.index[id]
	if !ok {
		return false
	}
	if !s.alive[idx] {
		return false
	}
	s.alive[idx] = false
	delete(s.index, id)
	s.aliveCount--
	return true
}

func (s *Store) MoveEntity(id uint64, dx, dy int32) bool {
	idx, ok := s.index[id]
	if !ok {
		return false
	}
	if !s.alive[idx] {
		return false
	}
	s.x[idx] += dx
	s.y[idx] += dy
	return true
}

func (s *Store) MoveEntityByIndex(idx int, dx, dy int32) bool {
	if idx < 0 || idx >= len(s.ids) {
		return false
	}
	if !s.alive[idx] {
		return false
	}
	s.x[idx] += dx
	s.y[idx] += dy
	return true
}

func (s *Store) EntityByID(id uint64) (Entity, bool) {
	idx, ok := s.index[id]
	if !ok {
		return Entity{}, false
	}
	if !s.alive[idx] {
		return Entity{}, false
	}
	return Entity{
		ID:   s.ids[idx],
		X:    s.x[idx],
		Y:    s.y[idx],
		Kind: s.kind[idx],
	}, true
}

func (s *Store) EntityIndex(id uint64) (int, bool) {
	idx, ok := s.index[id]
	if !ok {
		return 0, false
	}
	if !s.alive[idx] {
		return 0, false
	}
	return idx, true
}

func (s *Store) AppendEntities(dst []Entity) []Entity {
	for i, id := range s.ids {
		if !s.alive[i] {
			continue
		}
		dst = append(dst, Entity{
			ID:   id,
			X:    s.x[i],
			Y:    s.y[i],
			Kind: s.kind[i],
		})
	}
	return dst
}
