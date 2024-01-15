package events

type Type int

const (
	Unknown Type = iota
	Message
	Callback
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
