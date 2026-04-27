package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitProviderPath(t *testing.T) {
	tests := []struct {
		input      string
		wantPrefix string
		wantRest   string
		wantFound  bool
	}{
		{"local:/etc/hosts", "local", "/etc/hosts", true},
		{"my-mount:/some/path", "my-mount", "/some/path", true},
		{"/plain/path", "", "/plain/path", false},
		{"no-separator", "", "no-separator", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			prefix, rest, found := SplitProviderPath(tt.input)
			assert.Equal(t, tt.wantPrefix, prefix)
			assert.Equal(t, tt.wantRest, rest)
			assert.Equal(t, tt.wantFound, found)
		})
	}
}

func TestValidatePathSyntax(t *testing.T) {
	tests := []struct {
		path    string
		wantErr bool
	}{
		{"/", false},
		{"/valid/path", false},
		{"/path/with/slash/", true},
		{"/path/with/hash#", true},
		{"/path/with/question?", true},
		{"/path/with/star*", true},
		{"/path/with/bracket[", true},
		{"/path/with/backslash\\", true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			err := ValidatePathSyntax(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
