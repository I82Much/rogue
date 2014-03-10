package main

import (
	game "github.com/I82Much/rogue"
	termbox "github.com/nsf/termbox-go"
)

const (
	rows = 16
	cols = 32
)

func Render(w *game.World) {
	r := w.CurrentRoom()
	for row := 0; row < r.Rows(); row++ {
		for col := 0; col < r.Cols(); col++ {
			// col = x, row = y
			location := game.Loc(row, col)
			termbox.SetCell(col, row, r.RuneAt(location), termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.Flush()
	termbox.HideCursor()
}

func main() {

	// Set up controller
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.HideCursor()

	room1 := game.WalledRoom(rows, cols)
	room1.Spawn(rows/2, cols/2)
	room1.SpawnMonster()
	room1.SpawnMonster()

	// 2nd room to east of room 1
	room2 := game.WalledRoom(rows, cols)

	// 3rd room to south of room 1
	room3 := game.WalledRoom(rows, cols)

	world := game.NewWorld(2, 2)
	world.Set(game.Loc(0, 0), room1)
	world.Set(game.Loc(0, 1), room2)
	world.Set(game.Loc(1, 0), room3)

	// Set up doors between the rooms
	d1_2 := &game.Door{
		From: room1,
		To:   room2,
	}
	// Door to east
	room1.SetDoor(game.Loc(rows/2, cols-1), d1_2)
	d2_1 := &game.Door{
		From: room2,
		To:   room1,
		Same: d1_2,
	}
	d1_2.Same = d2_1
	room2.SetDoor(game.Loc(rows/2, 0), d2_1)

	d1_3 := &game.Door{
		From: room1,
		To:   room3,
	}
	// Door to south
	room1.SetDoor(game.Loc(rows-1, cols/2), d1_3)
	d3_1 := &game.Door{
		From: room3,
		To:   room1,
		Same: d1_3,
	}
	d1_3.Same = d3_1
	room3.SetDoor(game.Loc(0, cols/2), d3_1)

	Render(world)

	// Main game loop
	for {

		// Read input
		event := termbox.PollEvent()
		switch event.Key {
		case termbox.KeyArrowUp:
			world.MovePlayer(-1, 0)
		case termbox.KeyArrowRight:
			world.MovePlayer(0, 1)
		case termbox.KeyArrowDown:
			world.MovePlayer(1, 0)
		case termbox.KeyArrowLeft:
			world.MovePlayer(0, -1)
			// Quit
		case termbox.KeyCtrlC:
			return
		}
		// Render world
		Render(world)
	}
}
