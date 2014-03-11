package combat

import (
	"fmt"
	"sync"
	"time"

	termbox "github.com/nsf/termbox-go"
)

// The controller accepts input and converts it into commands for the model and view

type CombatController struct {
	Model   *CombatModel
	View    *CombatView
	runFlag bool

	unprocessedRunes []rune
	// Protects unprocessedRunes slice
	runesMutex sync.Mutex
}

func (c *CombatController) startup() {
	// Set up controller
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.HideCursor()

	// TODO(ndunn): This should be a channel not goroutine since we cannot stop this after it starts
	// Run the input polling loop in another goroutine
	go c.input()
}

func (c *CombatController) input() {
	for c.runFlag {
		event := termbox.PollEvent()
		if event.Key == termbox.KeyCtrlC {
			c.stop()
		}
		// They typed a letter
		if event.Ch != 0 {
			c.runesMutex.Lock()
			// TODO(ndunn): probably could/should use channels for this
			c.unprocessedRunes = append(c.unprocessedRunes, event.Ch)
			c.runesMutex.Unlock()
		}
	}
}

// drainUnprocessed returns all unprocessed runes, and sets the unprocessedRunes to nil.
func (c *CombatController) drainUnprocessed() []rune {
	c.runesMutex.Lock()
	defer c.runesMutex.Unlock()
	res := c.unprocessedRunes
	c.unprocessedRunes = nil
	return res
}

func (c *CombatController) update() {
	// Pull out the unprocessed runes
	c.Model.Update(c.drainUnprocessed())
}

func (c *CombatController) draw() {
	c.View.Render()
}

func shutdown() {
	termbox.Close()
	fmt.Printf("game over")
}

func (c *CombatController) stop() {
	c.runFlag = false
}

// Modified from http://entropyinteractive.com/2011/02/game-engine-design-the-game-loop/
func (c *CombatController) Run(dur time.Duration) {
	c.runFlag = true
	c.startup()
	nextTime := time.Now()
	for c.runFlag {
		now := time.Now()
		if now.Sub(nextTime) >= dur {
			nextTime = nextTime.Add(dur)
			c.update()
			// Really this should render some sort of victory/defeat screen
			/*if c.Model.Over() {
				c.stop()

			}*/
			c.draw()
		} else {
			sleepTime := nextTime.Sub(now)
			time.Sleep(sleepTime)
		}
	}
	shutdown()
}
