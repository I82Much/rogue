package dungeon

import (
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

func (c *Controller) Run() {
	c.running = true

	// Main game loop
	for c.running {
		// Read input
		event := termbox.PollEvent()
		switch event.Key {
		case termbox.KeyArrowUp:
			c.model.MovePlayer(-1, 0)

			/*if res := world.MovePlayer(-1, 0); res == game.CreatureOccupying {
				controller.Run(time.Duration(33) * time.Millisecond)
			}*/
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
