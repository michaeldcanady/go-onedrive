package abstractions

import "context"

type KeyValueStore interface {
	Get(context.Context, []byte) ([]byte, error)
	Set(context.Context, []byte, []byte) error
	Delete(context.Context, []byte) error
}
