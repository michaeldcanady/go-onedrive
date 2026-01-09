package event

import "errors"

var (
	ErrBusClosed            = errors.New("bus is closed")
	ErrUnknownTopic         = errors.New("unknown topic")
	ErrSubscriptionNotFound = errors.New("subscription not found")
)
