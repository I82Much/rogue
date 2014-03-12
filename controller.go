package rogue

import (
	"fmt"

	"github.com/I82Much/rogue/dungeon"
)

// Overarching controller for the entire game state
type Controller struct {
	model *Model
	view  *View
	running bool
	
	dungeonController *dungeon.Controller
}

func (c *Controller) Listen(e dungeon.Event) {
	if e == dungeon.EnterCombat {
		
	}
	
	fmt.Printf("Got event %v", e)
}

func NewController(model *Model, view *View) *Controller {
	c := &Controller{
		model: model,
		view:  view,
		running: false,
	}
	model.dungeonModel.AddListener(c)
	return c
}

func (c *Controller) Run() {
	c.running = true
	for c.running {
		
	}
}
