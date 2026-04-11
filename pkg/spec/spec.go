package spec

// Specification defines the interface for an object that can match a candidate of type T.
type Specification[T any] interface {
	IsSatisfiedBy(candidate T) bool
}

// AndSpecification combines two specifications using a logical AND.
type AndSpecification[T any] struct {
	left, right Specification[T]
}

func (s AndSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return s.left.IsSatisfiedBy(candidate) && s.right.IsSatisfiedBy(candidate)
}

// OrSpecification combines two specifications using a logical OR.
type OrSpecification[T any] struct {
	left, right Specification[T]
}

func (s OrSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return s.left.IsSatisfiedBy(candidate) || s.right.IsSatisfiedBy(candidate)
}

// NotSpecification inverts a specification using a logical NOT.
type NotSpecification[T any] struct {
	spec Specification[T]
}

func (s NotSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return !s.spec.IsSatisfiedBy(candidate)
}

// And returns a new specification that is the logical AND of two specifications.
func And[T any](left, right Specification[T]) Specification[T] {
	return AndSpecification[T]{left: left, right: right}
}

// Or returns a new specification that is the logical OR of two specifications.
func Or[T any](left, right Specification[T]) Specification[T] {
	return OrSpecification[T]{left: left, right: right}
}

// Not returns a new specification that is the logical NOT of a specification.
func Not[T any](spec Specification[T]) Specification[T] {
	return NotSpecification[T]{spec: spec}
}

// All returns a specification that always evaluates to true.
func All[T any]() Specification[T] {
	return AllSpecification[T]{}
}

type AllSpecification[T any] struct{}

func (s AllSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return true
}

// None returns a specification that always evaluates to false.
func None[T any]() Specification[T] {
	return NoneSpecification[T]{}
}

type NoneSpecification[T any] struct{}

func (s NoneSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return false
}
