package dungeon

import (
	"log"

	"github.com/I82Much/rogue/event"
	"github.com/I82Much/rogue/player"
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

func NewModule(w *World, p *player.Player) *Controller {
	model := NewModel(w, p)
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

func (c *Controller) HasWon() bool {
	return !c.model.world.AnyMonstersLeft()
}

func (c *Controller) MaybeUnlockCurrentRoom() {
	log.Println("maybe unlock room")
	curRoom := c.model.world.CurrentRoom()
	if !curRoom.AnyMonstersLeft() {
		log.Println("unlocking the doors")
		curRoom.UnlockDoors()
	} else {
		log.Println("monsters left; not unlocking the door")
	}
}

// This is a bit messy, but after successful combat we remember where we just fought
// (the tile we couldn't move onto because it was occupied), and then we remove the monster
// that was there and replace it with the player.
func (c *Controller) ReplaceMonsterWithPlayer() {
	// Bad.
	c.model.world.CurrentRoom().ReplaceMonsterWithPlayer(c.model.combatLocation)
}
