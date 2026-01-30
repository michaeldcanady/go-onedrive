package event

// Unsubscriber allows unsubscribing from events.
type Unsubscriber interface {
	// Unsubscribe unsubscribes a listener from a subscription ID.
	Unsubscribe(string, Listener) error
}
