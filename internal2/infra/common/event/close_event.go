package event

// CloseEvent represents an event indicating that the bus is closing.
type CloseEvent struct{}

// Topic returns the topic of the CloseEvent.
func (CloseEvent) Topic() string {
	return "bus/close"
}
