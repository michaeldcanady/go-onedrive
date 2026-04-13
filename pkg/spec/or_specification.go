package spec

// Or returns a new specification that is the logical OR of two specifications.
func Or[T any](left, right Specification[T]) Specification[T] {
	return OrSpecification[T]{left: left, right: right}
}

// OrSpecification combines two specifications using a logical OR.
type OrSpecification[T any] struct {
	left, right Specification[T]
}

func (s OrSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return s.left.IsSatisfiedBy(candidate) || s.right.IsSatisfiedBy(candidate)
}
