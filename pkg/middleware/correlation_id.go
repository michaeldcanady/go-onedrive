package middleware

import (
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

func WithCorrelationID(cmd *cobra.Command) {
	original := cmd.RunE

	if original == nil {
		return
	}

	cmd.RunE = func(c *cobra.Command, args []string) error {
		ctx := c.Context()

		cid := util.CorrelationIDFromContext(ctx)
		if cid == "" {
			cid = util.NewCorrelationID()
			ctx = util.WithCorrelationID(ctx, cid)
			c.SetContext(ctx)
		}

		return original(c, args)
	}
}
