package upload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    Options
		wantErr bool
	}{
		{
			name: "Valid options",
			opts: Options{
				Source:      "local.txt",
				Destination: "/remote.txt",
			},
			wantErr: false,
		},
		{
			name: "Missing source",
			opts: Options{
				Destination: "/remote.txt",
			},
			wantErr: true,
		},
		{
			name: "Missing destination",
			opts: Options{
				Source: "local.txt",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
