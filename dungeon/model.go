package dungeon

import (
	"github.com/I82Much/rogue/event"
	//	"fmt"
)

const (
	EnterCombat = "ENTER_COMBAT"
)

type Model struct {
	world     *World
	listeners []event.Listener

	combatLocation Location
}

func NewModel(w *World) *Model {
	return &Model{
		world: w,
	}
}

func (m *Model) AddListener(d event.Listener) {
	m.listeners = append(m.listeners, d)
}

func (m *Model) Publish(e string) {
	for _, listener := range m.listeners {
		listener.Listen(e)
	}
}

// TODO(ndunn): this is insane we now have 3 layers of wrappers. room -> world -> model. Do we even need this class? Probably not
func (m *Model) MovePlayer(rows, cols int) MovementResult {
	res := m.world.MovePlayer(rows, cols)
	if res == CreatureOccupying {
		// Store this for later
		m.combatLocation = m.world.CurrentRoom().playerLoc.Add(Loc(rows, cols))
		m.Publish(EnterCombat)
	}
	return res
}

// This is a bit messy, but after successful combat we remember where we just fought
// (the tile we couldn't move onto because it was occupied), and then we remove the monster
// that was there and replace it with the player.
func (m *Model) ReplaceMonsterWithPlayer() {
	m.world.CurrentRoom().ReplaceMonsterWithPlayer(m.combatLocation)
}
