package event

import (
	"container/list"
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
)

type InMemoryBus struct {
	lock      sync.RWMutex
	listeners map[string]*list.List
	closed    bool
	logger    logging.Logger
}

// NewInMemoryBus creates a new instance of InMemoryBus.
func NewInMemoryBus(logger logging.Logger) *InMemoryBus {
	return &InMemoryBus{
		listeners: make(map[string]*list.List),
		logger:    logger,
	}
}

// Close closes the event bus, preventing further subscriptions and publications.
// It notifies all listeners with a CloseEvent.
func (b *InMemoryBus) Close() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.closed {
		b.logger.Debug("event bus already closed")
		return nil
	}
	b.closed = true

	b.logger.Info("closing event bus; notifying listeners")

	for _, ll := range b.listeners {
		for e := ll.Front(); e != nil; e = e.Next() {
			sub := e.Value.(*subscription)
			if err := sub.listener.Listen(context.Background(), CloseEvent{}); err != nil {
				b.logger.Warn("listener returned error during CloseEvent",
					logging.String("subscription_id", sub.id),
					logging.Any("error", err),
				)
			}
		}
	}

	b.listeners = make(map[string]*list.List)
	return nil
}

// Publish publishes an event to its topic.
func (b *InMemoryBus) Publish(ctx context.Context, evt Topicer) error {
	b.lock.RLock()
	if b.closed {
		b.lock.RUnlock()
		b.logger.Error("publish attempted on closed bus", logging.String("topic", evt.Topic()))
		return ErrBusClosed
	}

	ll := b.listeners[evt.Topic()]
	b.lock.RUnlock()

	if ll == nil {
		b.logger.Debug("publish to topic with no listeners", logging.String("topic", evt.Topic()))
		return nil
	}

	for e := ll.Front(); e != nil; e = e.Next() {
		sub := e.Value.(*subscription)
		if err := sub.listener.Listen(ctx, evt); err != nil {
			b.logger.Error("listener returned error",
				logging.String("subscription_id", sub.id),
				logging.String("topic", evt.Topic()),
				logging.Any("error", err),
			)
			return err
		}
	}

	return nil
}

// Subscribe subscribes a listener to a topic and returns a subscription ID.
func (b *InMemoryBus) Subscribe(topic string, listener Listener) (string, error) {
	id := uuid.New().String()
	return id, b.SubscribeWithID(id, topic, listener)
}

// SubscribeWithID subscribes a listener to a topic with a specific subscription ID.
func (b *InMemoryBus) SubscribeWithID(id, topic string, listener Listener) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.closed {
		b.logger.Error("subscribe attempted on closed bus", logging.String("topic", topic))
		return ErrBusClosed
	}

	sub := newSubscription(id, listener)

	ll := b.listeners[topic]
	if ll == nil {
		ll = list.New()
		b.listeners[topic] = ll
	}

	ll.PushBack(sub)

	b.logger.Info("listener subscribed",
		logging.String("subscription_id", id),
		logging.String("topic", topic),
	)

	return nil
}

// Unsubscribe unsubscribes a listener from a subscription ID.
func (b *InMemoryBus) Unsubscribe(id string) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.closed {
		b.logger.Error("unsubscribe attempted on closed bus", logging.String("subscription_id", id))
		return ErrBusClosed
	}

	for topic, ll := range b.listeners {
		for e := ll.Front(); e != nil; e = e.Next() {
			sub := e.Value.(*subscription)
			if sub.id == id {
				ll.Remove(e)
				if ll.Len() == 0 {
					delete(b.listeners, topic)
				}

				b.logger.Info("listener unsubscribed",
					logging.String("subscription_id", id),
					logging.String("topic", topic),
				)

				return nil
			}
		}
	}

	b.logger.Warn("unsubscribe failed; subscription not found", logging.String("subscription_id", id))
	return ErrSubscriptionNotFound
}
