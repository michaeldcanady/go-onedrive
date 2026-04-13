package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrAll(t *testing.T) {
	tests := []struct {
		name     string
		specs    []Specification[any]
		expected bool
	}{
		{
			name:     "no specs",
			specs:    []Specification[any]{},
			expected: false,
		},
		{
			name:     "all false",
			specs:    []Specification[any]{mockSpec{false}, mockSpec{false}},
			expected: false,
		},
		{
			name:     "one true",
			specs:    []Specification[any]{mockSpec{false}, mockSpec{true}},
			expected: true,
		},
		{
			name:     "all true",
			specs:    []Specification[any]{mockSpec{true}, mockSpec{true}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := OrAll(tt.specs...)
			assert.Equal(t, tt.expected, s.IsSatisfiedBy("candidate"))
		})
	}
}

func TestOr(t *testing.T) {
	tests := []struct {
		name     string
		left     Specification[any]
		right    Specification[any]
		expected bool
	}{
		{
			name:     "no specs",
			left:     (Specification[any])(nil),
			right:    (Specification[any])(nil),
			expected: false,
		},
		{
			name:     "all false",
			left:     mockSpec{false},
			right:    mockSpec{false},
			expected: false,
		},
		{
			name:     "one true",
			left:     mockSpec{false},
			right:    mockSpec{true},
			expected: true,
		},
		{
			name:     "all true",
			left:     mockSpec{true},
			right:    mockSpec{true},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := OrAll(tt.left, tt.right)
			assert.Equal(t, tt.expected, s.IsSatisfiedBy("candidate"))
		})
	}
}
