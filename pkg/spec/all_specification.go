package spec

// All returns a specification that always evaluates to true.
func All[T any]() Specification[T] {
	return AllSpecification[T]{}
}

// AllSpecification is a specification that always evaluates to true.
type AllSpecification[T any] struct{}

func (s AllSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return true
}
