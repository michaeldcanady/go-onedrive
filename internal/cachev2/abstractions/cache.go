package abstractions

import "context"

type Cache[K comparable, V any] interface {
	NewEntry(context.Context, K) (*Entry[K, V], error)
	GetEntry(context.Context, K) (*Entry[K, V], error)
	SetEntry(context.Context, *Entry[K, V]) error
	Clear(context.Context) error
	Remove(K) error
	KeySerializer() Serializer[K]
}
