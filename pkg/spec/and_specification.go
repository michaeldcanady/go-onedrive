package spec

// And returns a new specification that is the logical AND of two specifications.
func And[T any](left, right Specification[T]) Specification[T] {
	return AndSpecification[T]{left: left, right: right}
}

// AndAll returns a new specification that is the logical AND of all provided specifications.
func AndAll[T any](specs ...Specification[T]) Specification[T] {
	return AndAllSpecification[T]{specs: specs}
}

// AndSpecification combines two specifications using a logical AND.
type AndSpecification[T any] struct {
	left, right Specification[T]
}

func (s AndSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return s.left.IsSatisfiedBy(candidate) && s.right.IsSatisfiedBy(candidate)
}

// AndAllSpecification combines multiple specifications using a logical AND.
type AndAllSpecification[T any] struct {
	specs []Specification[T]
}

func (s AndAllSpecification[T]) IsSatisfiedBy(candidate T) bool {
	result := true

	for _, spec := range s.specs {
		if spec == nil {
			continue
		}

		result = result && spec.IsSatisfiedBy(candidate)
	}
	return result
}
