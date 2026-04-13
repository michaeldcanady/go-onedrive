package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllSpecification(t *testing.T) {
	tests := []struct {
		name     string
		specs    []Specification[any]
		expected bool
	}{
		{
			name:     "no specs (always true)",
			specs:    []Specification[any]{},
			expected: true,
		},
		{
			name:     "all specs true",
			specs:    []Specification[any]{mockSpec{true}, mockSpec{true}},
			expected: true,
		},
		{
			name:     "one spec false",
			specs:    []Specification[any]{mockSpec{true}, mockSpec{false}},
			expected: false,
		},
		{
			name:     "all specs false",
			specs:    []Specification[any]{mockSpec{false}, mockSpec{false}},
			expected: false,
		},
		{
			name:     "with nil spec (others true)",
			specs:    []Specification[any]{mockSpec{true}, nil, mockSpec{true}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := All(tt.specs...)
			assert.Equal(t, tt.expected, s.IsSatisfiedBy("candidate"))
		})
	}
}
