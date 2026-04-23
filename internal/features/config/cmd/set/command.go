package set

import (
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/spf13/cobra"
)

func CreateSetCmd(container di.Container) *cobra.Command {
	var ctx *CommandContext
	opts := NewOptions()

	l, _ := container.Logger().CreateLogger("config-set")
	handler := NewCommand(container.Config(), l)

	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set configuration",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Key = args[0]
			opts.Value = args[1]
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			ctx = NewCommandContext(cmd.Context(), opts)
			if err := handler.Validate(ctx); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := handler.Execute(ctx); err != nil {
				return err
			}
			return handler.Finalize(ctx)
		},
	}

	return cmd
}
