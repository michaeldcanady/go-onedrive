package ls

import (
	"testing"

	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"
	"github.com/stretchr/testify/assert"
)

func TestOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		options Options
		wantErr bool
		errMsgs []string
	}{
		{
			name: "valid options",
			options: Options{
				Format: formatting.FormatShort,
			},
			wantErr: false,
		},
		{
			name: "invalid sort field",
			options: Options{
				Format:     formatting.FormatShort,
				SortFields: []string{"invalid"},
			},
			wantErr: true,
			errMsgs: []string{"invalid sorting field 'invalid'"},
		},
		{
			name: "unknown format",
			options: Options{
				Format: formatting.FormatUnknown,
			},
			wantErr: true,
			errMsgs: []string{"unknown output format specified"},
		},
		{
			name: "recursive incompatible format",
			options: Options{
				Recursive: true,
				Format:    formatting.FormatShort,
			},
			wantErr: true,
			errMsgs: []string{"recursive mode (-r/--recursive) is not supported with the 'short' format"},
		},
		{
			name: "recursive compatible format",
			options: Options{
				Recursive: true,
				Format:    formatting.FormatTree,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.options.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				for _, msg := range tt.errMsgs {
					assert.Contains(t, err.Error(), msg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
