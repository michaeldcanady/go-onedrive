package edit

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
			name:    "Valid path",
			opts:    Options{Path: "/folder/file.txt"},
			wantErr: false,
		},
		{
			name:    "Missing path",
			opts:    Options{Path: ""},
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
