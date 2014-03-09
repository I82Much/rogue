package rogue

// TODO implement world in terms of rooms

type World struct {
	rooms       [][]*Room
	currentRoom Location
}

func NewWorld(rows, cols int) *World {
	rooms := make([][]*Room, rows)
	for row := 0; row < rows; row++ {
		rooms[row] = make([]*Room, cols)
	}
	return &World{
		rooms:       rooms,
		currentRoom: Loc(0, 0),
	}
}

func (w *World) RoomAt(loc Location) *Room {
	return w.rooms[loc.Row][loc.Col]
}

func (w *World) Set(loc Location, r *Room) {
	w.rooms[loc.Row][loc.Col] = r
}

func (w *World) MovePlayer(rows, cols int) MovementResult {
	r := w.CurrentRoom()
	return r.MovePlayer(rows, cols)
}

func (w *World) CurrentRoom() *Room {
	return w.RoomAt(w.currentRoom)
}
