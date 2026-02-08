package options

type Option[T any] = func(T) error
