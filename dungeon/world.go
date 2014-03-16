package dungeon

import (
	"fmt"
)

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

func (w *World) AnyMonstersLeft() bool {
	for row := 0; row < len(w.rooms); row++ {
		for col := 0; col < len(w.rooms[row]); col++ {
			if w.rooms[row][col].AnyMonstersLeft() {
				return true
			}
		}
	}
	return false
}

func (w *World) RoomAt(loc Location) *Room {
	return w.rooms[loc.Row][loc.Col]
}

func (w *World) Set(loc Location, r *Room) {
	w.rooms[loc.Row][loc.Col] = r
}

func (w *World) MovePlayer(rows, cols int) MovementResult {
	r := w.CurrentRoom()
	res := r.MovePlayer(rows, cols)
	// Did player move onto a door
	if res == Move && r.PlayerTile() == UnlockedDoor {
		d := r.doors[r.playerLoc]
		if d.From == nil || d.To == nil {
			panic(fmt.Sprintf("nil from/to for door at loc %v", r.playerLoc))
		}

		// Delete player from old room
		r.RemovePlayer()
		// spawn him into new room
		newRoom := d.To

		linkedDoor := d.Same
		doorLoc := newRoom.DoorLocation(linkedDoor)
		if doorLoc == nil {
			panic(fmt.Sprintf("no linked door for %v", *linkedDoor))
		}
		// Where did user end up after walking through door
		destination := doorLoc.Add(Loc(rows, cols))
		newRoom.Spawn(destination.Row, destination.Col)
		// Find which of the rooms is the linked door
		loc := InvalidLoc
		for row := 0; row < len(w.rooms); row++ {
			for col := 0; col < len(w.rooms[0]); col++ {
				if w.RoomAt(Loc(row, col)) == newRoom {
					loc = Loc(row, col)
				}
			}
		}
		if loc == InvalidLoc {
			panic(fmt.Sprintf("couldn't find room in world whose value is %v", *newRoom))
		}
		w.currentRoom = loc
	}
	return res
}

func (w *World) CurrentRoom() *Room {
	return w.RoomAt(w.currentRoom)
}
