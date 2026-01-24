package di

import "github.com/michaeldcanady/go-onedrive/internal/event"

type EventBus interface {
	event.Publisher
	event.Subscriber
	event.IDSubscriber
}
