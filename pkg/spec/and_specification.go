package spec

// And returns a new specification that is the logical AND of two specifications.
func And[T any](left, right Specification[T]) Specification[T] {
	return AndSpecification[T]{left: left, right: right}
}

// AndSpecification combines two specifications using a logical AND.
type AndSpecification[T any] struct {
	left, right Specification[T]
}

func (s AndSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return s.left.IsSatisfiedBy(candidate) && s.right.IsSatisfiedBy(candidate)
}
