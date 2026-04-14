package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type person struct {
	name string
	age  int
}

func TestSorter(t *testing.T) {
	items := []person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 30},
		{"Dave", 20},
	}

	compareByName := func(i, j person) bool { return i.name < j.name }
	compareByAge := func(i, j person) bool { return i.age < j.age }

	t.Run("Sort by age", func(t *testing.T) {
		s := NewSorter(items)
		s.AddComparator(compareByAge)
		sorted := s.Sort()

		expected := []person{
			{"Dave", 20},
			{"Bob", 25},
			{"Alice", 30},
			{"Charlie", 30},
		}
		assert.Equal(t, expected, sorted)
		// Ensure original slice is unchanged
		assert.Equal(t, "Alice", items[0].name)
	})

	t.Run("Sort by age then name", func(t *testing.T) {
		s := NewSorter(items)
		s.AddComparator(compareByAge)
		s.AddComparator(compareByName)
		sorted := s.Sort()

		expected := []person{
			{"Dave", 20},
			{"Bob", 25},
			{"Alice", 30},
			{"Charlie", 30},
		}
		assert.Equal(t, expected, sorted)
	})

	t.Run("Sort by name", func(t *testing.T) {
		s := NewSorter(items)
		s.AddComparator(compareByName)
		sorted := s.Sort()

		expected := []person{
			{"Alice", 30},
			{"Bob", 25},
			{"Charlie", 30},
			{"Dave", 20},
		}
		assert.Equal(t, expected, sorted)
	})
}
