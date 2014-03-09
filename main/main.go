package main

import (
	"fmt"

	game "github.com/I82Much/roguelike"
)

func main() {
	world := game.NewWorld(32, 32)
	fmt.Printf("%v", world)
}
