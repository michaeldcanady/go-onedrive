package spec

// None returns a specification that always evaluates to false.
func None[T any](specs ...Specification[T]) Specification[T] {
	return NoneSpecification[T]{specs: specs}
}

// NoneSpecification is a specification that always evaluates to false.
type NoneSpecification[T any] struct {
	specs []Specification[T]
}

func (s NoneSpecification[T]) IsSatisfiedBy(candidate T) bool {
	for _, spec := range s.specs {
		if spec == nil {
			continue
		}
		// If any specification is satisfied, then NoneSpecification should return false.
		if spec.IsSatisfiedBy(candidate) {
			return false
		}
	}
	return true
}
