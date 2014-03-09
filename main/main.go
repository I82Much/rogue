package main

import (
	"fmt"

	game "github.com/I82Much/roguelike"
	termbox "github.com/nsf/termbox-go"
)

const (
	rows = 16
	cols = 32
)

func main() {
	world := game.NewWorld(rows, cols)
	// Make the outline set to WALL
	for i := 0; i < cols; i++ {
		// Top row
		world.Set(game.Loc(0, i), game.Wall)
		// Bottom row
		world.Set(game.Loc(rows-1, i), game.Wall)
	}
	for i := 0; i < rows; i++ {
		// Top row
		world.Set(game.Loc(i, 0), game.Wall)
		// Bottom row
		world.Set(game.Loc(i, cols-1), game.Wall)
	}
	world.Spawn(rows/2, cols/2)

	fmt.Printf("%v\n", world)

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
			world.MovePlayer(-1, 0)
		}

		// Render world
		fmt.Printf("%v\n", world)

	}
}
