package main

import (
	"fmt"

	game "github.com/I82Much/roguelike"
)

const (
	rows = 32
	cols = 32
)

func main() {
	world := game.NewWorld(rows, cols)
	world.Set(rows/2, cols/2, game.PLAYER)
	// Make the outline set to WALL
	for i := 0; i < cols; i++ {
		// Top row
		world.Set(0, i, game.WALL)
		// Bottom row
		world.Set(rows-1, i, game.WALL)
	}
	for i := 0; i < rows; i++ {
		// Top row
		world.Set(i, 0, game.WALL)
		// Bottom row
		world.Set(i, cols-1, game.WALL)
	}

	fmt.Printf("%v", world)
}
