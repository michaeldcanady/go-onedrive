package spec

// NotSpecification inverts a specification using a logical NOT.
type NotSpecification[T any] struct {
	spec Specification[T]
}

func (s NotSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return !s.spec.IsSatisfiedBy(candidate)
}

// Not returns a new specification that is the logical NOT of a specification.
func Not[T any](spec Specification[T]) Specification[T] {
	return NotSpecification[T]{spec: spec}
}
