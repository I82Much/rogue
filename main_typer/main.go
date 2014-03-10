package main

import (
	"fmt"
	"time"

	termbox "github.com/nsf/termbox-go"
)

/**
 * @param value The incoming value to be converted
 * @param low1  Lower bound of the value's current range
 * @param high1 Upper bound of the value's current range
 * @param low2  Lower bound of the value's target range
 * @param high2 Upper bound of the value's target range
 */
func doMap(value, low1, high1, low2, high2 float64) float64 {
	diff := value - low1
	proportion := diff / (high1 - low1)
	return lerp(low2, high2, proportion)
}

// Linearly interpolate between two values
func lerp(value1, value2, amt float64) float64 {
	return ((value2 - value1) * amt) + value1
}


type Challenge struct {
	word string
	// How much has been spelled so far
	spelled string
	start time.Time
	hits, attempts int
}

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

	hits := 0
	attempts := 0

	// There are still words to spell
	for len(spelledWords) < len(words) {
		for _, toSpell := range words {
			start := time.Now()
			maxDuration := time.Duration(10) * time.Second

			runes := []rune(toSpell)
			spelled := ""
			// Main game loop

			// Still more to render
			for len(spelled) < len(toSpell) {
				elapsed := time.Now().Sub(start)
				if elapsed > maxDuration {
					fmt.Println("you're too slow")
				}

				maxRows := 20
				curRow := int(doMap(elapsed.Seconds(), 0.0, maxDuration.Seconds(), 0, float64(maxRows-1)))

				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

				// Render to spell on screen
				for i, s := range toSpell {
					spelled := i < len(spelled)
					foregroundColor := termbox.ColorDefault
					if spelled {
						foregroundColor = termbox.ColorRed
					}
					termbox.SetCell(i, curRow, s, foregroundColor, termbox.ColorDefault)
				}
				if attempts > 0 {
					accuracyText := fmt.Sprintf("%d / %d (%.2f%%)", hits, attempts, 100.0*float32(hits)/float32(attempts))
					for i, c := range accuracyText {
						termbox.SetCell(i, curRow+3, c, termbox.ColorDefault, termbox.ColorDefault)
					}
				}
				termbox.Flush()

				// Read input

				event := termbox.PollEvent()
				if event.Ch == runes[len(spelled)] {
					spelled = spelled + string(event.Ch)
					hits++
					attempts++
				} else {
					attempts++
				}
				if event.Key == termbox.KeyCtrlC {
					return
				}
			}

			spelledWords = append(spelledWords, toSpell)
		}
	}
}
