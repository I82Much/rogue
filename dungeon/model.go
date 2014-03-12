package dungeon

type Event int
const (
	EnterCombat Event = iota
)

type Model struct {
	world *World
	// Where we emit events to
	e chan<- Event
}

func NewModel(w *World, e chan<- Event) *Model {
	return &Model{
		world: w,
		e: e,
	}
}

// TODO(ndunn): this is insane we now have 3 layers of wrappers. room -> world -> model. Do we even need this class? Probably not
func (m *Model) MovePlayer(rows, cols int) MovementResult {
	res := m.world.MovePlayer(rows, cols)
	if res == CreatureOccupying {
		go func() {
			m.e <- EnterCombat
		}()
	}
	return res
}
