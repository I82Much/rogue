package dungeon

import (
	"fmt"
	"math/rand"

	"github.com/I82Much/rogue/monster"
)

type World struct {
	rooms       [][]*Room
	currentRoom Location
}

const (
	minMonstersPerRoom = 1
	maxMonstersPerRoom = 5

	// 20% of the time we'll stack two consecutive monsters on the same spot
	probOfStacking = float32(0.2)
)

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

func RandomWorld(rows, cols int) *World {
	w := NewWorld(rows, cols)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			doorLocs := make(map[DoorDir]bool)
			// TODO(ndunn): it'd be better if not every room were completely connected, but this is an easier algorithm.
			if row != 0 {
				doorLocs[North] = true
			}
			if row != rows-1 {
				doorLocs[South] = true
			}
			if col != 0 {
				doorLocs[West] = true
			}
			if col != cols-1 {
				doorLocs[East] = true
			}

			room := RandomRoom(doorLocs)
			// TODO(ndunn): automatically add all the doors between different rooms
			w.Set(Loc(row, col), room)
		}
	}

	// Construct the portals between the rooms. This is a pain in the butt.
	for row := 0; row < cols; row++ {
		for col := 0; col < cols; col++ {
			room := w.rooms[row][col]

			var northRoom, eastRoom, southRoom, westRoom *Room
			if row != 0 {
				northRoom = w.rooms[row-1][col]
			}
			if row != rows-1 {
				southRoom = w.rooms[row+1][col]
			}
			if col != 0 {
				westRoom = w.rooms[row][col-1]
			}
			if col != cols-1 {
				eastRoom = w.rooms[row][col+1]
			}

			// Set the doors between them
			if northRoom != nil {
				maybeJoinAdjacentNorthSouthRooms(northRoom, room)
			}
			if southRoom != nil {
				maybeJoinAdjacentNorthSouthRooms(room, southRoom)
			}
			if westRoom != nil {
				maybeJoinAdjacentWestEastRooms(westRoom, room)
			}
			if eastRoom != nil {
				maybeJoinAdjacentWestEastRooms(room, eastRoom)
			}

			// Only fill the room after setting the doors so that the monsters don't show up right next to door
			fillRoomWithMonsters(room)
		}
	}

	// Spawn player in a random room
	randRow := int(rand.Int31n(int32(rows)))
	randCol := int(rand.Int31n(int32(cols)))
	w.currentRoom = Loc(randRow, randCol)
	w.CurrentRoom().RandSpawn()
	return w
}

func maybeJoinAdjacentNorthSouthRooms(north, south *Room) {
	d1_2 := &Door{
		From: north,
		To:   south,
	}
	d2_1 := &Door{
		From: south,
		To:   north,
	}
	d1_2.Same = d2_1
	d2_1.Same = d1_2

	north.SetDoor(South, d1_2)
	south.SetDoor(North, d2_1)
}

func maybeJoinAdjacentWestEastRooms(west, east *Room) {
	d1_2 := &Door{
		From: west,
		To:   east,
	}
	d2_1 := &Door{
		From: east,
		To:   west,
	}
	d1_2.Same = d2_1
	d2_1.Same = d1_2

	west.SetDoor(East, d1_2)
	east.SetDoor(West, d2_1)
}

func fillRoomWithMonsters(r *Room) {
	lastLoc := InvalidLoc
	for _, m := range randMonsters() {
		// some % of the time, we stack monsters
		if lastLoc != InvalidLoc && rand.Float32() < probOfStacking {
			r.AddMonster(lastLoc, m)
		} else {
			lastLoc, _ = r.SpawnMonster(m)
		}
	}
}

func randMonsters() []monster.Type {
	numMonsters := minMonstersPerRoom
	if minMonstersPerRoom != maxMonstersPerRoom {
		 numMonsters = minMonstersPerRoom + int(rand.Int31n(int32(maxMonstersPerRoom-minMonstersPerRoom)))
	}
	var types []monster.Type
	for i := 0; i < numMonsters; i++ {
		monsterIndex := rand.Perm(len(monster.All))[0]
		types = append(types, monster.All[monsterIndex])
	}
	if len(types) < minMonstersPerRoom {
		panic(fmt.Sprintf("expected at least %d monsters; got %d", minMonstersPerRoom, len(types)))
	}
	return types
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
