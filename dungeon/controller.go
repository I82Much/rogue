package dungeon

import (
	"github.com/I82Much/rogue/event"
	termbox "github.com/nsf/termbox-go"
)

type Controller struct {
	model   *Model
	view    *View
	running bool
}

func NewController(model *Model, view *View) *Controller {
	return &Controller{
		model:   model,
		view:    view,
		running: false,
	}
}

func NewModule(w *World) *Controller {
	model := NewModel(w)
	view := NewView(model)
	return NewController(model, view)
}

func (c *Controller) AddListener(d event.Listener) {
	c.model.AddListener(d)
}

func (c *Controller) Start() {
	c.running = true
	c.view.Render()
	// Main game loop
	for c.running {
		// Read input
		event := termbox.PollEvent()
		switch event.Key {
		case termbox.KeyArrowUp:
			c.model.MovePlayer(-1, 0)
		case termbox.KeyArrowRight:
			c.model.MovePlayer(0, 1)
		case termbox.KeyArrowDown:
			c.model.MovePlayer(1, 0)
		case termbox.KeyArrowLeft:
			c.model.MovePlayer(0, -1)
			// Quit
		case termbox.KeyCtrlC:
			return
		}
		// Render world
		c.view.Render()
	}
}

func (c *Controller) Stop() {
	c.running = false
}
