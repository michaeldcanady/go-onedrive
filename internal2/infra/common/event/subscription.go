package event

type subscription struct {
	id       string
	listener Listener
}

func newSubscription(id string, listener Listener) *subscription {
	return &subscription{
		id:       id,
		listener: listener,
	}
}
