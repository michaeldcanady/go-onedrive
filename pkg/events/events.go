package events

import (
	"sync"
)

// Event is the interface that all events must implement.
type Event interface {
	Name() string
}

// Handler is a function that processes an event.
type Handler func(Event)

// Dispatcher manages the registration and execution of event handlers.
type Dispatcher struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

// NewDispatcher initializes a new Event Dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string][]Handler),
	}
}

// Subscribe registers a handler for a specific event name.
func (d *Dispatcher) Subscribe(eventName string, handler Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eventName] = append(d.handlers[eventName], handler)
}

// Dispatch broadcasts an event to all registered handlers.
func (d *Dispatcher) Dispatch(event Event) {
	d.mu.RLock()
	handlers, ok := d.handlers[event.Name()]
	d.mu.RUnlock()

	if !ok {
		return
	}

	for _, handler := range handlers {
		handler(event)
	}
}
