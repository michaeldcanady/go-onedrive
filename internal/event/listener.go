package event

import "context"

type Listener interface {
	Listen(context.Context, Topicer) error
}

type ListenerFunc func(ctx context.Context, evt Topicer) error

func (f ListenerFunc) Listen(ctx context.Context, evt Topicer) error {
	return f(ctx, evt)
}
