package ls

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
			name: "Valid defaults",
			opts: Options{
				Format:       "short",
				SortProperty: "name",
			},
			wantErr: false,
		},
		{
			name: "Conflict: FoldersOnly and FilesOnly",
			opts: Options{
				FoldersOnly: true,
				FilesOnly:   true,
			},
			wantErr: true,
		},
		{
			name: "Invalid Format",
			opts: Options{
				Format: "invalid",
			},
			wantErr: true,
		},
		{
			name: "Invalid SortProperty",
			opts: Options{
				Format:       "short",
				SortProperty: "invalid",
			},
			wantErr: true,
		},
		{
			name: "All supported formats",
			opts: Options{
				Format:       "json",
				SortProperty: "name",
			},
			wantErr: false,
		},
		{
			name: "Supported property: size",
			opts: Options{
				Format:       "short",
				SortProperty: "size",
			},
			wantErr: false,
		},
		{
			name: "Supported property: modified",
			opts: Options{
				Format:       "short",
				SortProperty: "modified",
			},
			wantErr: false,
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
