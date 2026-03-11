package formatting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatterFactory_Create(t *testing.T) {
	factory := NewFormatterFactory()

	tests := []struct {
		format   string
		expected interface{}
		wantErr  bool
	}{
		{"short", &HumanShortFormatter{}, false},
		{"", &HumanShortFormatter{}, false},
		{"long", &HumanLongFormatter{}, false},
		{"json", &JSONFormatter{}, false},
		{"yaml", &YAMLFormatter{}, false},
		{"yml", &YAMLFormatter{}, false},
		{"tree", &TreeFormatter{}, false},
		{"invalid", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			result, err := factory.Create(tt.format)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, tt.expected, result)
			}
		})
	}
}
