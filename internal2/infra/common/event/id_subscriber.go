package event

// IDSubscriber allows subscribing with a specific subscription ID.
type IDSubscriber interface {
	// SubscribeWithID subscribes with a specific subscription ID.
	SubscribeWithID(string, string, Listener) error
}
