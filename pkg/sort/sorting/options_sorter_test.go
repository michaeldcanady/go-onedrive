package sorting

import (
	"testing"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/stretchr/testify/assert"
)

func TestOptionsSorter(t *testing.T) {
	now := time.Now()
	items := []fs.Item{
		{Name: "b", Size: 20, ModifiedAt: now.Add(-1 * time.Hour)},
		{Name: "a", Size: 30, ModifiedAt: now.Add(-2 * time.Hour)},
		{Name: "a", Size: 10, ModifiedAt: now},
		{Name: "c", Size: 10, ModifiedAt: now},
	}

	t.Run("Sort by Name Ascending", func(t *testing.T) {
		opts := SortingOptions{
			Criteria: []SortingCriteria{
				{Field: "Name", Direction: DirectionAscending},
			},
		}
		sorter := NewOptionsSorter(opts)
		sorted, err := sorter.Sort(items)
		assert.NoError(t, err)
		assert.Equal(t, "a", sorted[0].Name)
		assert.Equal(t, "a", sorted[1].Name)
		assert.Equal(t, "b", sorted[2].Name)
		assert.Equal(t, "c", sorted[3].Name)
	})

	t.Run("Sort by Name then Size", func(t *testing.T) {
		opts := SortingOptions{
			Criteria: []SortingCriteria{
				{Field: "Name", Direction: DirectionAscending},
				{Field: "Size", Direction: DirectionAscending},
			},
		}
		sorter := NewOptionsSorter(opts)
		sorted, err := sorter.Sort(items)
		assert.NoError(t, err)
		assert.Equal(t, "a", sorted[0].Name)
		assert.Equal(t, int64(10), sorted[0].Size)
		assert.Equal(t, "a", sorted[1].Name)
		assert.Equal(t, int64(30), sorted[1].Size)
	})

	t.Run("Sort by Size then Name", func(t *testing.T) {
		opts := SortingOptions{
			Criteria: []SortingCriteria{
				{Field: "Size", Direction: DirectionAscending},
				{Field: "Name", Direction: DirectionAscending},
			},
		}
		sorter := NewOptionsSorter(opts)
		sorted, err := sorter.Sort(items)
		assert.NoError(t, err)
		assert.Equal(t, int64(10), sorted[0].Size)
		assert.Equal(t, "a", sorted[0].Name)
		assert.Equal(t, int64(10), sorted[1].Size)
		assert.Equal(t, "c", sorted[1].Name)
	})

	t.Run("Unknown Field", func(t *testing.T) {
		opts := SortingOptions{
			Criteria: []SortingCriteria{
				{Field: "Unknown", Direction: DirectionAscending},
			},
		}
		sorter := NewOptionsSorter(opts)
		_, err := sorter.Sort(items)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown sort field")
	})
}
