package app

type Option[T any] = func(T) error
