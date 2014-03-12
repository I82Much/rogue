package event

type Listener interface {
	Listen(s string)
}
