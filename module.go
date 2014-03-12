package rogue

import (
	"github.com/I82Much/rogue/event"
)

type Module interface {
	AddListener(listener event.Listener)
	Start()
	Stop()
}