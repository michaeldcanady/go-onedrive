package event

import "context"

type Listener interface {
	Listen(context.Context, Topicer) error
}
