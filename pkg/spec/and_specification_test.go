package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnd(t *testing.T) {
	tests := []struct {
		name     string
		left     Specification[any]
		right    Specification[any]
		expected bool
	}{
		{
			name:     "both true",
			left:     mockSpec{true},
			right:    mockSpec{true},
			expected: true,
		},
		{
			name:     "left false",
			left:     mockSpec{false},
			right:    mockSpec{true},
			expected: false,
		},
		{
			name:     "right false",
			left:     mockSpec{true},
			right:    mockSpec{false},
			expected: false,
		},
		{
			name:     "both false",
			left:     mockSpec{false},
			right:    mockSpec{false},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := And(tt.left, tt.right)
			assert.Equal(t, tt.expected, s.IsSatisfiedBy("candidate"))
		})
	}
}

func TestAndAll(t *testing.T) {
	tests := []struct {
		name     string
		specs    []Specification[any]
		expected bool
	}{
		{
			name:     "no specs",
			specs:    []Specification[any]{},
			expected: true,
		},
		{
			name:     "all true",
			specs:    []Specification[any]{mockSpec{true}, mockSpec{true}, mockSpec{true}},
			expected: true,
		},
		{
			name:     "one false in middle",
			specs:    []Specification[any]{mockSpec{true}, mockSpec{false}, mockSpec{true}},
			expected: false,
		},
		{
			name:     "all false",
			specs:    []Specification[any]{mockSpec{false}, mockSpec{false}, mockSpec{false}},
			expected: false,
		},
		{
			name:     "with nil spec (all others true)",
			specs:    []Specification[any]{mockSpec{true}, nil, mockSpec{true}},
			expected: true,
		},
		{
			name:     "with nil spec (one false)",
			specs:    []Specification[any]{mockSpec{true}, nil, mockSpec{false}},
			expected: false,
		},
		{
			name:     "nil-only specs",
			specs:    []Specification[any]{nil, nil},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := AndAll(tt.specs...)
			assert.Equal(t, tt.expected, s.IsSatisfiedBy("candidate"))
		})
	}
}
