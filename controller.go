package rogue

import (
	"fmt"

	"github.com/I82Much/rogue/dungeon"
)

// Overarching controller for the entire game state
type Controller struct {
	model *Model
	view  *View
}

func (c *Controller) Listen(e dungeon.Event) {
	fmt.Printf("Got event %v", e)
}

func NewController(model *Model, view *View) *Controller {
	c := &Controller{
		model: model,
		view:  view,
	}
	model.dungeonModel.AddListener(c)
	return c
}

func (c *Controller) Run() {

}
