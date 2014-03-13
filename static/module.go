// Helper functions for creating modules from static screens
package static

import (
	"github.com/I82Much/rogue/event"
	"github.com/I82Much/rogue/render"
	termbox "github.com/nsf/termbox-go"
)

type Module struct {
	contents string
	// The event that should be published if the given rune is pressed
	keyMap    map[rune]string
	listeners []event.Listener
	running   bool
}

func (s *Module) AddListener(d event.Listener) {
	s.listeners = append(s.listeners, d)
}

func NewModule(contents string, keyMap map[rune]string) *Module {
	return &Module{
		contents: contents,
		keyMap:   keyMap,
	}
}

func (s *Module) Render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	render.Render(s.contents, 0, 0)
	termbox.Flush()
}

func (s *Module) Start() {
	s.running = true
	s.Render()
	s.input()
}

func (s *Module) publish(evt string) {
	for _, d := range s.listeners {
		d.Listen(evt)
	}
}

func (s *Module) Stop() {
	s.running = false
}

func (s *Module) input() {
	for s.running {
		event := termbox.PollEvent()
		// See if we have any events to publish based on this event.
		pubEvent, ok := s.keyMap[event.Ch]
		if !ok {
			continue
		}
		s.publish(pubEvent)
	}
}
