package main

import (
	"time"

	termbox "github.com/nsf/termbox-go"
)

var (
	runFlag = false
	maxRows = 21
	words  = []*Word{
		&Word{
			word: "Hello",
			row: 0,
			onScreen: time.Now(),
			duration: time.Duration(3) * time.Second,
		},
		&Word{
			word: "World",
			row: 2,
			onScreen: time.Now(),
			duration: time.Duration(4) * time.Second,
		},
		&Word{
			word: "Supercalifragilistic",
			row: 4,
			onScreen: time.Now(),
			duration: time.Duration(2) * time.Second,
		},
	}
)

type Word struct {
	word string
	row int
	onScreen time.Time 
	duration time.Duration
}



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

func main() {
	run(time.Duration(33) * time.Millisecond)
}

func startup() {
	// Set up controller
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.HideCursor()
}

func update() {
	now := time.Now()
	anyOnscreen := false
	for _, word := range words {
		elapsed := now.Sub(word.onScreen)
		row := int(doMap(elapsed.Seconds(), 0.0, word.duration.Seconds(), 0, float64(maxRows-1)))
		word.row = row
		if row < maxRows {
			anyOnscreen = true
		}
	}
	// If we've run out of words, end it.
	if !anyOnscreen {
		stop()
	}
}

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, word := range words {
		for i, c := range word.word {
			termbox.SetCell(i, word.row, c, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.Flush()
}

func shutdown() {
	termbox.Close()
}

func stop() {
	runFlag = false
}

// Modified from http://entropyinteractive.com/2011/02/game-engine-design-the-game-loop/
func run(dur time.Duration) {
	runFlag = true
	startup()
	nextTime := time.Now()
	for runFlag {
		now := time.Now()
		if now.Sub(nextTime) >= dur {
			//fmt.Printf("rendering at %v", now)
			nextTime = nextTime.Add(dur)
			update()
			draw()
		} else {
			sleepTime := nextTime.Sub(now)
			time.Sleep(sleepTime)
		}
	}
	shutdown()
}
