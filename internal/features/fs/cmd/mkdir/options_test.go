package mkdir

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		options Options
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid path",
			options: Options{
				Path: "/some/path",
			},
			wantErr: false,
		},
		{
			name: "empty path",
			options: Options{
				Path: "",
			},
			wantErr: true,
			errMsg:  "path is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.options.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
