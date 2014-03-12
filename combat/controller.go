package combat

import (
	"fmt"
	"sync"
	"time"

	"github.com/I82Much/rogue/event"
	termbox "github.com/nsf/termbox-go"
)

// The controller accepts input and converts it into commands for the model and view

type Controller struct {
	Model   *Model
	View    *View
	runFlag bool
	unprocessedRunes []rune
	// Protects unprocessedRunes slice
	runesMutex sync.Mutex
}

func NewModule(player *Player, monsters []*Monster) *Controller {
	model := NewCombatModel(player, monsters)
	view := &View{
		Model: model,
		rows: 20,
	}
	return &Controller {
		Model: model,
		View: view,
	}
}

func (c *Controller) AddListener(d event.Listener) {
	c.Model.AddListener(d)
}

func (c *Controller) startup() {
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

func (c *Controller) input() {
	for c.runFlag {
		event := termbox.PollEvent()
		// TODO(ndunn): This should be done somewhere else.
		if event.Key == termbox.KeyCtrlC {
			c.Stop()
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
func (c *Controller) drainUnprocessed() []rune {
	c.runesMutex.Lock()
	defer c.runesMutex.Unlock()
	res := c.unprocessedRunes
	c.unprocessedRunes = nil
	return res
}

func (c *Controller) update() {
	// Pull out the unprocessed runes
	c.Model.Update(c.drainUnprocessed())
}

func (c *Controller) draw() {
	c.View.Render()
}

func shutdown() {
	termbox.Close()
	fmt.Printf("game over")
}

func (c *Controller) Start() {
	c.startup()
	c.Run(time.Duration(33) * time.Millisecond)
}

func (c *Controller) Stop() {
	c.runFlag = false
}

// Modified from http://entropyinteractive.com/2011/02/game-engine-design-the-game-loop/
func (c *Controller) Run(dur time.Duration) {
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
