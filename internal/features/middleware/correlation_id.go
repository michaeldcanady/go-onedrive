package middleware

import (
	"github.com/michaeldcanady/go-onedrive/internal/features/shared"
	"github.com/spf13/cobra"
)

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
