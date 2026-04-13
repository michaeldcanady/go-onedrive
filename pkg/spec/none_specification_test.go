package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoneSpecification(t *testing.T) {
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
			name:     "all false",
			specs:    []Specification[any]{mockSpec{false}, mockSpec{false}},
			expected: true,
		},
		{
			name:     "one true",
			specs:    []Specification[any]{mockSpec{false}, mockSpec{true}},
			expected: false,
		},
		{
			name:     "all true",
			specs:    []Specification[any]{mockSpec{true}, mockSpec{true}},
			expected: false,
		},
		{
			name:     "with nil spec (all others false)",
			specs:    []Specification[any]{mockSpec{false}, nil, mockSpec{false}},
			expected: true,
		},
		{
			name:     "with nil spec (one true)",
			specs:    []Specification[any]{mockSpec{false}, nil, mockSpec{true}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := None(tt.specs...)
			assert.Equal(t, tt.expected, s.IsSatisfiedBy("candidate"))
		})
	}
}
