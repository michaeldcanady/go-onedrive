package event

import "context"

// Publisher allows publishing events.
type Publisher interface {
	// Publish publishes an event to its topic.
	Publish(context.Context, Topicer) error
}
