package dungeon

import (
//	"fmt"
)

type Event string

const (
	EnterCombat Event = "ENTER_COMBAT"
)

type DungeonListener interface {
	Listen(e Event)
}

type Model struct {
	world     *World
	listeners []DungeonListener
}

func NewModel(w *World) *Model {
	return &Model{
		world: w,
	}
}

func (m *Model) AddListener(d DungeonListener) {
	m.listeners = append(m.listeners, d)
}

func (m *Model) Publish(e Event) {
	//	fmt.Printf("publishing %v\n", e)
	for _, listener := range m.listeners {
		//		fmt.Printf("publishing %v\n", e)
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
