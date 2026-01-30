package event

// Subscriber allows subscribing to events.
type Subscriber interface {
	// Subscribe subscribes to a topic with a listener.
	Subscribe(string, Listener) (string, error)
}
