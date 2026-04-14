package spec

// Specification defines the interface for an object that can match a candidate of type T.
type Specification[T any] interface {
	IsSatisfiedBy(candidate T) bool
}
