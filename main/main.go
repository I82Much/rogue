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

	// Door to east
	room1.SetTile(game.Loc(rows/2, cols-1), game.Door)
	// Door to south
	room1.SetTile(game.Loc(rows-1, cols/2), game.Door)

	// 2nd room to east of room 1
	room2 := game.WalledRoom(rows, cols)
	// Door to west back to room 1
	room2.SetTile(game.Loc(rows/2, 0), game.Door)

	// 3rd room to south of room 1
	room3 := game.WalledRoom(rows, cols)
	// Door to north back to room 1
	room3.SetTile(game.Loc(0, cols/2), game.Door)

	world := game.NewWorld(2, 2)
	world.Set(game.Loc(0, 0), room1)
	world.Set(game.Loc(0, 1), room2)
	world.Set(game.Loc(1, 0), room3)
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
