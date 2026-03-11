package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    ItemType
		expected string
	}{
		{"file", ItemTypeFile, "file"},
		{"folder", ItemTypeFolder, "folder"},
		{"unknown", ItemTypeUnknown, "unknown"},
		{"invalid enum value", ItemType(999), "unknown"},
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
		expected ItemType
	}{
		{"file lowercase", "file", ItemTypeFile},
		{"file uppercase", "FILE", ItemTypeFile},
		{"folder lowercase", "folder", ItemTypeFolder},
		{"folder uppercase", "FOLDER", ItemTypeFolder},
		{"unknown lowercase", "unknown", ItemTypeUnknown},
		{"unknown uppercase", "UNKNOWN", ItemTypeUnknown},
		{"invalid string", "not-a-type", ItemTypeUnknown},
		{"empty string", "", ItemTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, ParseItemType(tt.input))
		})
	}
}

func TestItemType_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    ItemType
		expected string
	}{
		{"file", ItemTypeFile, `"file"`},
		{"folder", ItemTypeFolder, `"folder"`},
		{"unknown", ItemTypeUnknown, `"unknown"`},
		{"invalid enum value", ItemType(999), `"unknown"`},
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
		expected    ItemType
		expectError bool
	}{
		{"file", `"file"`, ItemTypeFile, false},
		{"folder", `"folder"`, ItemTypeFolder, false},
		{"unknown", `"unknown"`, ItemTypeUnknown, false},
		{"invalid string", `"not-a-type"`, ItemTypeUnknown, false},
		{"empty string", `""`, ItemTypeUnknown, false},
		{"invalid JSON type (number)", `123`, ItemTypeUnknown, true},
		{"invalid JSON type (object)", `{}`, ItemTypeUnknown, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var it ItemType
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
		input ItemType
	}{
		{"file", ItemTypeFile},
		{"folder", ItemTypeFolder},
		{"unknown", ItemTypeUnknown},
		{"invalid enum value", ItemType(999)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b, err := json.Marshal(tt.input)
			require.NoError(t, err)

			var out ItemType
			err = json.Unmarshal(b, &out)
			require.NoError(t, err)

			// invalid enum values normalize to unknown
			if tt.input != ItemTypeFile && tt.input != ItemTypeFolder {
				assert.Equal(t, ItemTypeUnknown, out)
			} else {
				assert.Equal(t, tt.input, out)
			}
		})
	}
}
