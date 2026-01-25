package cliprofileservicego

type Option[T any] = func(T) error
