package filtering

import (
	"testing"
	"time"

	shared "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/stretchr/testify/assert"
)

func TestFilterOptions_Apply(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		opts     []FilterOption
		validate func(*testing.T, *FilterOptions)
	}{
		{
			name: "WithItemType",
			opts: []FilterOption{WithItemType(shared.TypeFile)},
			validate: func(t *testing.T, o *FilterOptions) {
				assert.Equal(t, shared.TypeFile, o.ItemType)
			},
		},
		{
			name: "IncludeAll",
			opts: []FilterOption{IncludeAll()},
			validate: func(t *testing.T, o *FilterOptions) {
				assert.True(t, o.IncludeAll)
			},
		},
		{
			name: "ExcludeHidden",
			opts: []FilterOption{ExcludeHidden()},
			validate: func(t *testing.T, o *FilterOptions) {
				assert.False(t, o.IncludeAll)
			},
		},
		{
			name: "WithName",
			opts: []FilterOption{WithName("*.txt"), WithName("*.jpg")},
			validate: func(t *testing.T, o *FilterOptions) {
				assert.Equal(t, []string{"*.txt", "*.jpg"}, o.Names)
			},
		},
		{
			name: "WithMinSize",
			opts: []FilterOption{WithMinSize(100)},
			validate: func(t *testing.T, o *FilterOptions) {
				assert.Equal(t, int64(100), *o.MinSize)
			},
		},
		{
			name: "WithMaxSize",
			opts: []FilterOption{WithMaxSize(200)},
			validate: func(t *testing.T, o *FilterOptions) {
				assert.Equal(t, int64(200), *o.MaxSize)
			},
		},
		{
			name: "WithModifiedBefore",
			opts: []FilterOption{WithModifiedBefore(now)},
			validate: func(t *testing.T, o *FilterOptions) {
				assert.True(t, now.Equal(*o.ModifiedBefore))
			},
		},
		{
			name: "WithModifiedAfter",
			opts: []FilterOption{WithModifiedAfter(now)},
			validate: func(t *testing.T, o *FilterOptions) {
				assert.True(t, now.Equal(*o.ModifiedAfter))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewFilterOptions()
			err := o.Apply(tt.opts)
			assert.NoError(t, err)
			tt.validate(t, o)
		})
	}
}
