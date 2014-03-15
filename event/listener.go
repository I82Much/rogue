package event

type Listener interface {
	// Extra data can be passed in the 'extras' piece. A bit of a hack.
	Listen(event string, extras interface{})
}
