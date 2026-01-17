package abstractions

type Entry[K, V any] struct {
	key   K
	value V
}

func NewEntry[K, V any](key K, value V) *Entry[K, V] {
	return &Entry[K, V]{
		key:   key,
		value: value,
	}
}

// GetKey implements [abstraction.Entry].
func (c *Entry[K, V]) GetKey() K {
	return c.key
}

// GetValue implements [abstraction.Entry].
func (c *Entry[K, V]) GetValue() V {
	return c.value
}

// SetValue implements [abstraction.Entry].
func (c *Entry[K, V]) SetValue(value V) {
	c.value = value
}
