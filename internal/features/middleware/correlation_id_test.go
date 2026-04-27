package middleware

import (
	"context"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/core/shared"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestWithCorrelationID(t *testing.T) {
	tests := []struct {
		name           string
		initialContext func() context.Context
		wantErr        bool
	}{
		{
			name: "injects new correlation id if missing",
			initialContext: func() context.Context {
				return context.Background()
			},
			wantErr: false,
		},
		{
			name: "preserves existing correlation id",
			initialContext: func() context.Context {
				return shared.WithCorrelationID(context.Background(), "existing-id")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executed := false
			cmd := &cobra.Command{
				RunE: func(c *cobra.Command, args []string) error {
					executed = true
					ctx := c.Context()
					cid := shared.CorrelationIDFromContext(ctx)
					assert.NotEmpty(t, cid)

					if tt.name == "preserves existing correlation id" {
						assert.Equal(t, "existing-id", cid)
					}
					return nil
				},
			}

			cmd.SetContext(tt.initialContext())
			WithCorrelationID(cmd)

			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.True(t, executed)
		})
	}
}

func TestWithCorrelationID_NilRunE(t *testing.T) {
	cmd := &cobra.Command{}
	WithCorrelationID(cmd)
	assert.Nil(t, cmd.RunE)
}
