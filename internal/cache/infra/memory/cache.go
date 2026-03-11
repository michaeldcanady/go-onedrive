package infra

import (
	"github.com/michaeldcanady/go-onedrive/pkg/cache/memory"
)

type Cache[K comparable, V any] = memory.Cache[K, V]

func New[K comparable, V any]() *Cache[K, V] {
	return memory.New[K, V]()
}
