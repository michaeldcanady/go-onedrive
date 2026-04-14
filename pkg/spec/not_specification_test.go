package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotSpecification(t *testing.T) {
	tests := []struct {
		name     string
		spec     Specification[any]
		expected bool
	}{
		{
			name:     "not true is false",
			spec:     mockSpec{result: true},
			expected: false,
		},
		{
			name:     "not false is true",
			spec:     mockSpec{result: false},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Not(tt.spec)
			assert.Equal(t, tt.expected, s.IsSatisfiedBy("candidate"))
		})
	}
}
