package spec

// None returns a specification that always evaluates to false.
func None[T any]() Specification[T] {
	return NoneSpecification[T]{}
}

// NoneSpecification is a specification that always evaluates to false.
type NoneSpecification[T any] struct{}

func (s NoneSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return false
}
