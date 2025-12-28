package cache

import "errors"

var (
	ErrAcquireLock    = errors.New("unable to acquire lock for key")
	ErrCreateCacheDir = errors.New("unable to create cache directory")
)
