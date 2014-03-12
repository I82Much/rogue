package dungeon

type Model struct {
	world *World
}

func NewModel(w *World) *Model {
	return &Model{
		world: w,
	}
}

// TODO(ndunn): this is insane we now have 3 layers of wrappers. room -> world -> model. Do we even need this class? Probably not
func (m *Model) MovePlayer(rows, cols int) MovementResult {
	return m.world.MovePlayer(rows, cols)
}
