package main

import (
	"fmt"
	"sync"
	"time"

	termbox "github.com/nsf/termbox-go"
)

var (
	runFlag = false
	maxRows = 21
	words   = []*Word{
		&Word{
			word:     "Hello",
			row:      0,
			onScreen: time.Now(),
			duration: time.Duration(3) * time.Second,
		},
		&Word{
			word:     "World",
			row:      2,
			onScreen: time.Now(),
			duration: time.Duration(4) * time.Second,
		},
		&Word{
			word:     "Supercalifragilistic",
			row:      4,
			onScreen: time.Now(),
			duration: time.Duration(10) * time.Second,
		},
	}

	currentTyping  *Word
	hits, attempts int

	completedWords int

	unprocessedRunes = []rune{}
	// Protects unprocessedRunes slice
	runesMutex sync.Mutex
)

type Word struct {
	word     string
	spelled  []rune
	row      int
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
	// Run the input polling loop in another goroutine
	go input()
}

func input() {
	for runFlag {
		event := termbox.PollEvent()
		if event.Key == termbox.KeyCtrlC {
			stop()
		}
		// They typed a letter
		if event.Ch != 0 {
			runesMutex.Lock()
			// TODO(ndunn): probably could/should use channels for this
			unprocessedRunes = append(unprocessedRunes, event.Ch)
			runesMutex.Unlock()
		}
	}
}

func removeWord(w *Word) {
	index := -1
	for i, word := range words {
		if word == w {
			words = append(words[0:i], words[i+1:]...)
			return
		}
	}
	if index == -1 {
		panic(fmt.Sprintf("couldn't find word %v", w))
	}
}

func update() {
	now := time.Now()

	// Pull out the unprocessed runes, apply them to the words
	runesMutex.Lock()
	for _, r := range unprocessedRunes {
		attempts++
		// Does this rune represent the first untyped letter of any of the candidates? If so it's a hit. If not it's a miss
		if currentTyping != nil {
			runes := []rune(currentTyping.word)
			if r == runes[len(currentTyping.spelled)] {
				hits++
				currentTyping.spelled = append(currentTyping.spelled, r)

				// Done the word
				if len(currentTyping.spelled) == len(currentTyping.word) {
					removeWord(currentTyping)
					completedWords++
					currentTyping = nil
				}
			}
		} else {
			for _, word := range words {
				runes := []rune(word.word)
				if r == runes[len(word.spelled)] {
					hits++
					word.spelled = append(word.spelled, r)
					currentTyping = word
					break
				}
			}
		}
	}
	// Clear out unprocessed runes - they've been processed
	unprocessedRunes = nil
	runesMutex.Unlock()

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
		foreground := termbox.ColorDefault
		if word == currentTyping {
			foreground = foreground | termbox.AttrBold
		}
		for i, c := range word.word {
			if i < len(word.spelled) {
				termbox.SetCell(i, word.row, c, termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)
			} else {
				termbox.SetCell(i, word.row, c, foreground, termbox.ColorDefault)
			}
		}
	}

	runesMutex.Lock()
	for i, r := range unprocessedRunes {
		termbox.SetCell(i, 0, r, termbox.ColorRed, termbox.ColorDefault)
	}
	runesMutex.Unlock()

	if attempts > 0 {
		accuracyText := fmt.Sprintf("%d / %d (%.2f%%)", hits, attempts, 100.0*float32(hits)/float32(attempts))
		for i, c := range accuracyText {
			termbox.SetCell(i, 3, c, termbox.ColorDefault, termbox.ColorDefault)
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
