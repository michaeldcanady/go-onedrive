package infra

import (
	"github.com/michaeldcanady/go-onedrive/pkg/cache"
	"github.com/michaeldcanady/go-onedrive/pkg/cache/bolt"
)

type Store = bolt.Store

var (
	NewStore        = bolt.NewStore
	NewSiblingStore = bolt.NewSiblingStore
	NewStoreWithDB  = bolt.NewStoreWithDB
	ErrKeyNotFound  = cache.ErrKeyNotFound
)
