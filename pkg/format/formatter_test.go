package format

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatters_StructuredData(t *testing.T) {
	data := map[string]any{
		"client_id":     "test-id",
		"client_secret": "test-secret",
	}

	factory := NewFactory()

	t.Run("JSON Formatter", func(t *testing.T) {
		f := factory.Get(FormatJSON)
		var buf bytes.Buffer
		err := f.Format(&buf, data)
		assert.NoError(t, err)

		expected := "{\n  \"client_id\": \"test-id\",\n  \"client_secret\": \"test-secret\"\n}\n"
		assert.Equal(t, expected, buf.String())
	})

	t.Run("YAML Formatter", func(t *testing.T) {
		f := factory.Get(FormatYAML)
		var buf bytes.Buffer
		err := f.Format(&buf, data)
		assert.NoError(t, err)

		// YAML output can vary slightly in indentation/newlines, but should be structured
		assert.Contains(t, buf.String(), "client_id: test-id")
		assert.Contains(t, buf.String(), "client_secret: test-secret")
	})
}

func BenchmarkFormatters(b *testing.B) {
	data := map[string]any{
		"client_id":     "test-id",
		"client_secret": "test-secret",
		"complex_nested": map[string]any{
			"key1": "value1",
			"key2": 123,
			"key3": []string{"a", "b", "c"},
		},
	}

	factory := NewFactory()

	b.Run("JSON", func(b *testing.B) {
		f := factory.Get(FormatJSON)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = f.Format(io.Discard, data)
		}
	})

	b.Run("YAML", func(b *testing.B) {
		f := factory.Get(FormatYAML)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = f.Format(io.Discard, data)
		}
	})
}
