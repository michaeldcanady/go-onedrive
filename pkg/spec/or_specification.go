package spec

// Or returns a new specification that is the logical OR of two specifications.
func Or[T any](left, right Specification[T]) Specification[T] {
	return OrSpecification[T]{left: left, right: right}
}

// OrAll returns a new specification that is the logical OR of all provided specifications.
func OrAll[T any](specs ...Specification[T]) Specification[T] {
	return OrAllSpecification[T]{specs: specs}
}

// OrSpecification combines two specifications using a logical OR.
type OrSpecification[T any] struct {
	left, right Specification[T]
}

func (s OrSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return s.left.IsSatisfiedBy(candidate) || s.right.IsSatisfiedBy(candidate)
}

// OrAllSpecification combines multiple specifications using a logical OR.
type OrAllSpecification[T any] struct {
	specs []Specification[T]
}

func (s OrAllSpecification[T]) IsSatisfiedBy(candidate T) bool {
	for _, spec := range s.specs {
		if spec == nil {
			continue
		}
		if spec.IsSatisfiedBy(candidate) {
			return true
		}
	}
	return false
}
