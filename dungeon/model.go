package dungeon

import (
	"github.com/I82Much/rogue/event"
	"github.com/I82Much/rogue/player"
)

const (
	EnterCombat = "ENTER_COMBAT"
)

type Model struct {
	world     *World
	player *player.Player
	listeners []event.Listener

	combatLocation Location
}

func NewModel(w *World, p *player.Player) *Model {
	return &Model{
		world: w,
		player: p,
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

