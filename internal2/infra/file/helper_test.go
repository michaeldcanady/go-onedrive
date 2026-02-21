package file

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty", "", ""},
		{"slash", "/", ""},
		{"dot", ".", ""},
		{"simple", "Documents", "/Documents"},
		{"leading slash", "/Documents", "/Documents"},
		{"trailing slash", "Documents/", "/Documents"},
		{"nested", "foo/bar", "/foo/bar"},
		{"double slashes", "//foo//bar//", "/foo/bar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.expected, normalizePath(tt.input))
		})
	}
}

func TestDeref(t *testing.T) {
	t.Parallel()

	type sampleStruct struct {
		A int
		B string
	}

	tests := []struct {
		name     string
		input    any
		ptr      any
		expected any
	}{
		{
			name:     "nil int pointer returns zero",
			ptr:      (*int)(nil),
			expected: 0,
		},
		{
			name: "non-nil int pointer returns value",
			ptr: func() *int {
				v := 42
				return &v
			}(),
			expected: 42,
		},
		{
			name:     "nil string pointer returns empty string",
			ptr:      (*string)(nil),
			expected: "",
		},
		{
			name: "non-nil string pointer returns value",
			ptr: func() *string {
				v := "hello"
				return &v
			}(),
			expected: "hello",
		},
		{
			name:     "nil struct pointer returns zero struct",
			ptr:      (*sampleStruct)(nil),
			expected: sampleStruct{},
		},
		{
			name: "non-nil struct pointer returns struct value",
			ptr: func() *sampleStruct {
				v := sampleStruct{A: 10, B: "x"}
				return &v
			}(),
			expected: sampleStruct{A: 10, B: "x"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch ptr := tt.ptr.(type) {

			case *int:
				require.Equal(t, tt.expected, deref(ptr))

			case *string:
				require.Equal(t, tt.expected, deref(ptr))

			case *sampleStruct:
				require.Equal(t, tt.expected, deref(ptr))

			default:
				t.Fatalf("unsupported test type: %T", tt.ptr)
			}
		})
	}
}
