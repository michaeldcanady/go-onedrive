package abstractions

import "context"

type MetadataProvider[K comparable, M any] interface {
	GetMetadata(context.Context, K) (M, error)
	SetMetadata(context.Context, K, M) error
}
