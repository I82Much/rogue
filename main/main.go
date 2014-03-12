package main

import (
	"github.com/I82Much/rogue"
	termbox "github.com/nsf/termbox-go"
)

func main() {
	// Set up controller
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.HideCursor()
	defer termbox.Close()

	game := rogue.NewGame()
	game.Start()
}
