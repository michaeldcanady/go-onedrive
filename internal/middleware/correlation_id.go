package middleware

import (
	"github.com/michaeldcanady/go-onedrive/internal/shared"
	"github.com/spf13/cobra"
)

// WithCorrelationID is a middleware that ensures a correlation ID is present in the command's context. If not, it generates a new one and adds it to the context.
func WithCorrelationID(cmd *cobra.Command) {
	original := cmd.RunE

	if original == nil {
		return
	}

	cmd.RunE = func(c *cobra.Command, args []string) error {
		ctx := c.Context()

		cid := shared.CorrelationIDFromContext(ctx)
		if cid == "" {
			cid = shared.NewCorrelationID()
			ctx = shared.WithCorrelationID(ctx, cid)
			c.SetContext(ctx)
		}

		return original(c, args)
	}
}
