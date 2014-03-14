package main

import (
	"fmt"
	"log"
	"os"

	"github.com/I82Much/rogue"
	termbox "github.com/nsf/termbox-go"
)

func main() {

	f, err := os.OpenFile("testlogfile2", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer f.Close()
	log.SetOutput(f)

	// Set up controller
	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.HideCursor()
	defer termbox.Close()

	game := rogue.NewGame()
	log.Println("starting game")
	game.Start()
}
