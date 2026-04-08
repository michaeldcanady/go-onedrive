package middleware

import (
	"context"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/shared"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestWithCorrelationID(t *testing.T) {
	tests := []struct {
		name           string
		initialContext func() context.Context
		checkCID       func(t *testing.T, cid string)
	}{
		{
			name: "Correlation ID is set in context",
			initialContext: func() context.Context {
				return context.Background()
			},
			checkCID: func(t *testing.T, cid string) {
				assert.NotEmpty(t, cid, "Correlation ID should be set in context")
			},
		},
		{
			name: "Existing Correlation ID is preserved",
			initialContext: func() context.Context {
				return shared.WithCorrelationID(context.Background(), "existing-correlation-id")
			},
			checkCID: func(t *testing.T, cid string) {
				assert.Equal(t, "existing-correlation-id", cid, "Existing Correlation ID should be preserved")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{
				RunE: func(c *cobra.Command, args []string) error {
					ctx := c.Context()
					cid := shared.CorrelationIDFromContext(ctx)
					tt.checkCID(t, cid)
					return nil
				},
			}

			cmd.SetContext(tt.initialContext())

			WithCorrelationID(cmd)

			err := cmd.RunE(cmd, []string{})
			assert.NoError(t, err)
		})
	}
}
