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

	room := game.NewRoom(rows, cols)
	// Make the outline set to WALL
	for i := 0; i < cols; i++ {
		// Top row
		room.SetTile(game.Loc(0, i), game.Wall)
		// Bottom row
		room.SetTile(game.Loc(rows-1, i), game.Wall)
	}
	for i := 0; i < rows; i++ {
		// Top row
		room.SetTile(game.Loc(i, 0), game.Wall)
		// Bottom row
		room.SetTile(game.Loc(i, cols-1), game.Wall)
	}
	room.Spawn(rows/2, cols/2)

	room.SpawnMonster()
	room.SpawnMonster()
	
	// Door to north
	room.SetTile(game.Loc(0, cols/2), game.Door)
	// Door to east
	room.SetTile(game.Loc(rows/2, cols-1), game.Door)
	// Door to west
	room.SetTile(game.Loc(rows/2, 0), game.Door)
	// Door to south
	room.SetTile(game.Loc(rows-1, cols/2), game.Door)

	world := game.NewWorld(1, 1)
	world.Set(game.Loc(0, 0), room)
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
