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
		m.Publish(EnterCombat)
	}
	return res
}
