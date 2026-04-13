package spec

// All returns a specification that always evaluates to true.
func All[T any](specs ...Specification[T]) Specification[T] {
	return AllSpecification[T]{
		specs: specs,
	}
}

// AllSpecification is a specification that always evaluates to true.
type AllSpecification[T any] struct {
	specs []Specification[T]
}

func (s AllSpecification[T]) IsSatisfiedBy(candidate T) bool {
	for _, spec := range s.specs {
		if spec == nil {
			continue
		}
		// If any specification is not satisfied, then AllSpecification should return false.
		if !spec.IsSatisfiedBy(candidate) {
			return false
		}
	}
	return true
}
