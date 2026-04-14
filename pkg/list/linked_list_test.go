package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedList(t *testing.T) {
	t.Run("Empty list", func(t *testing.T) {
		l := New[int]()
		assert.Equal(t, 0, l.Len())
		assert.Nil(t, l.Front())
		assert.Nil(t, l.Back())
	})

	t.Run("PushBack", func(t *testing.T) {
		l := New[int]()
		l.PushBack(1)
		l.PushBack(2)
		l.PushBack(3)

		assert.Equal(t, 3, l.Len())
		assert.Equal(t, 1, l.Front().Value())
		assert.Equal(t, 3, l.Back().Value())

		// Iterate
		var values []int
		for e := l.Front(); e != nil; e = e.Next() {
			values = append(values, e.Value())
		}
		assert.Equal(t, []int{1, 2, 3}, values)
	})

	t.Run("PushFront", func(t *testing.T) {
		l := New[int]()
		l.PushFront(1)
		l.PushFront(2)
		l.PushFront(3)

		assert.Equal(t, 3, l.Len())
		assert.Equal(t, 3, l.Front().Value())
		assert.Equal(t, 1, l.Back().Value())

		// Iterate
		var values []int
		for e := l.Front(); e != nil; e = e.Next() {
			values = append(values, e.Value())
		}
		assert.Equal(t, []int{3, 2, 1}, values)
	})

	t.Run("Prev", func(t *testing.T) {
		l := New[int]()
		l.PushBack(1)
		l.PushBack(2)
		l.PushBack(3)

		e := l.Back()
		assert.Equal(t, 3, e.Value())
		e = e.Prev()
		assert.Equal(t, 2, e.Value())
		e = e.Prev()
		assert.Equal(t, 1, e.Value())
		e = e.Prev()
		assert.Nil(t, e)
	})
}
