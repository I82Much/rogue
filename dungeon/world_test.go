package dungeon

import (
	"testing"
)

func TestDoors(t *testing.T) {
	rows := 16
	cols := 32

	room1 := WalledRoom(rows, cols)
	// 2nd room to east of room 1
	room2 := WalledRoom(rows, cols)
	world := NewWorld(2, 2)
	world.Set(Loc(0, 0), room1)
	world.Set(Loc(0, 1), room2)

	// Set up doors between the rooms
	d1_2 := &Door{
		From: room1,
		To:   room2,
	}
	// Door to east
	room1.SetDoor(Loc(rows/2, cols-1), d1_2)
	d2_1 := &Door{
		From: room2,
		To:   room1,
		Same: d1_2,
	}
	d1_2.Same = d2_1
	room2.SetDoor(Loc(rows/2, 0), d2_1)

	// Spawn player next to door
	room1.Spawn(rows/2, cols-2)
	// Move him east through door
	world.MovePlayer(0, 1)
	// Player should be in room 2 now, 1 space east of that room's door
	if world.CurrentRoom() != room2 {
		t.Errorf("expected player to have moved through door into room 2; was in room %v", *world.CurrentRoom())
	}
	if room2.playerLoc != Loc(rows/2, 1) {
		t.Errorf("expected player to be one tile east of east room's door; was at loc %v", room2.playerLoc)
	}

	// Move him west through door - should end up back in room 1
	world.MovePlayer(0, -1)
	if world.CurrentRoom() != room1 {
		t.Errorf("expected player to have moved through door back into room 1; was in room %v", *world.CurrentRoom())
	}
	if room1.playerLoc != Loc(rows/2, cols-2) {
		t.Errorf("expected player to be one tile west of west room's door; was at loc %v", room1.playerLoc)
	}
}
