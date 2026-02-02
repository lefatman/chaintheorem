// File: internal/world/tick.go
package world

type Tick struct {
	Seq uint32
}

func (t *Tick) Step(store *Store, intents *IntentStore) uint32 {
	t.Seq++
	ApplyIntents(store, intents)
	return t.Seq
}

func ApplyIntents(store *Store, intents *IntentStore) int {
	if store == nil || intents == nil {
		return 0
	}
	applied := 0
	for i := range intents.playerIDs {
		if !intents.active[i] {
			continue
		}
		intent := intents.intents[i]
		intents.active[i] = false
		if intent.DX == 0 && intent.DY == 0 {
			continue
		}
		entityIndex := intents.entityIndex[i]
		if entityIndex == invalidEntityIndex {
			continue
		}
		if store.MoveEntityByIndex(entityIndex, intent.DX, intent.DY) {
			applied++
		}
	}
	return applied
}
