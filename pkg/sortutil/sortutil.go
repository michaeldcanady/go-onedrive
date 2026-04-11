package sortutil

// Comparator defines the interface for an object that can compare two candidates of type T.
type Comparator[T any] interface {
	// Less returns true if i is considered strictly less than j.
	Less(i, j T) bool
}

// ComparatorFunc is an adapter to allow the use of ordinary functions as comparators.
type ComparatorFunc[T any] func(i, j T) bool

// Less calls f(i, j).
func (f ComparatorFunc[T]) Less(i, j T) bool {
	return f(i, j)
}

// Then combines two comparators to create a primary and secondary ordering.
// If the primary comparator cannot distinguish between i and j (neither i < j nor j < i),
// the secondary comparator is used.
func Then[T any](primary, secondary Comparator[T]) Comparator[T] {
	return ComparatorFunc[T](func(i, j T) bool {
		if primary.Less(i, j) {
			return true
		}
		if primary.Less(j, i) {
			return false
		}
		return secondary.Less(i, j)
	})
}

// Reverse inverts the logic of the provided comparator.
func Reverse[T any](comparator Comparator[T]) Comparator[T] {
	return ComparatorFunc[T](func(i, j T) bool {
		return comparator.Less(j, i)
	})
}

// AllAlwaysGreater returns a comparator where every item is "greater" than everything else (useful for nil cases).
type AllAlwaysGreater[T any] struct{}

func (s AllAlwaysGreater[T]) Less(i, j T) bool {
	return false
}
