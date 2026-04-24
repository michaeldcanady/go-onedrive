package validation

// Policy defines the interface for an object that can validate a candidate of type T.
type Policy[T any] interface {
	Evaluate(candidate T) error
}

// PolicyFunc is an adapter to allow the use of ordinary functions as validation policies.
type PolicyFunc[T any] func(candidate T) error

// Evaluate calls f(candidate).
func (f PolicyFunc[T]) Evaluate(candidate T) error {
	return f(candidate)
}
