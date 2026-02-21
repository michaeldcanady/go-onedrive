package file_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
)

func TestItemType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    file.ItemType
		expected string
	}{
		{"file", file.ItemTypeFile, "file"},
		{"folder", file.ItemTypeFolder, "folder"},
		{"unknown", file.ItemTypeUnknown, "unknown"},
		{"invalid enum value", file.ItemType(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, tt.input.String())
		})
	}
}

func TestParseItemType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected file.ItemType
	}{
		{"file lowercase", "file", file.ItemTypeFile},
		{"file uppercase", "FILE", file.ItemTypeFile},
		{"folder lowercase", "folder", file.ItemTypeFolder},
		{"folder uppercase", "FOLDER", file.ItemTypeFolder},
		{"unknown lowercase", "unknown", file.ItemTypeUnknown},
		{"unknown uppercase", "UNKNOWN", file.ItemTypeUnknown},
		{"invalid string", "not-a-type", file.ItemTypeUnknown},
		{"empty string", "", file.ItemTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, file.ParseItemType(tt.input))
		})
	}
}

func TestItemType_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    file.ItemType
		expected string
	}{
		{"file", file.ItemTypeFile, `"file"`},
		{"folder", file.ItemTypeFolder, `"folder"`},
		{"unknown", file.ItemTypeUnknown, `"unknown"`},
		{"invalid enum value", file.ItemType(999), `"unknown"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b, err := json.Marshal(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(b))
		})
	}
}

func TestItemType_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       string
		expected    file.ItemType
		expectError bool
	}{
		{"file", `"file"`, file.ItemTypeFile, false},
		{"folder", `"folder"`, file.ItemTypeFolder, false},
		{"unknown", `"unknown"`, file.ItemTypeUnknown, false},
		{"invalid string", `"not-a-type"`, file.ItemTypeUnknown, false},
		{"empty string", `""`, file.ItemTypeUnknown, false},
		{"invalid JSON type (number)", `123`, file.ItemTypeUnknown, true},
		{"invalid JSON type (object)", `{}`, file.ItemTypeUnknown, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var it file.ItemType
			err := json.Unmarshal([]byte(tt.input), &it)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, it)
		})
	}
}

func TestItemType_JSONRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input file.ItemType
	}{
		{"file", file.ItemTypeFile},
		{"folder", file.ItemTypeFolder},
		{"unknown", file.ItemTypeUnknown},
		{"invalid enum value", file.ItemType(999)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b, err := json.Marshal(tt.input)
			require.NoError(t, err)

			var out file.ItemType
			err = json.Unmarshal(b, &out)
			require.NoError(t, err)

			// invalid enum values normalize to unknown
			if tt.input != file.ItemTypeFile && tt.input != file.ItemTypeFolder {
				assert.Equal(t, file.ItemTypeUnknown, out)
			} else {
				assert.Equal(t, tt.input, out)
			}
		})
	}
}
