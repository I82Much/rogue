package main

import (
	termbox "github.com/nsf/termbox-go"
)

func main() {
	// Set up controller
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.HideCursor()

	words := []string{
		"hello",
		"world",
		"supercalifragilisticexpialidocious",
	}
	spelledWords := []string{}
	// There are still words to spell
	for len(spelledWords) < len(words) {
		for _, toSpell := range words {
			runes := []rune(toSpell)
			spelled := ""
			// Main game loop
			for len(spelled) < len(toSpell) {
				// Render to spell on screen
				for i, s := range toSpell {
					spelled := i < len(spelled)
					foregroundColor := termbox.ColorDefault
					if spelled {
						foregroundColor = termbox.ColorRed
					}
					termbox.SetCell(i, 0, s, foregroundColor, termbox.ColorDefault)
				}
				termbox.Flush()

				// Read input

				event := termbox.PollEvent()
				if event.Ch == runes[len(spelled)] {
					spelled = spelled + string(event.Ch)
				}
				if event.Key == termbox.KeyCtrlC {
					return
				}
			}

			spelledWords = append(spelledWords, toSpell)
		}
	}
}
