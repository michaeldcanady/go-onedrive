package get

import (
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/spf13/cobra"
)

func CreateGetCmd(container di.Container) *cobra.Command {
	var ctx *CommandContext
	opts := NewOptions()

	l, _ := container.Logger().CreateLogger("config-get")
	handler := NewCommand(container.Config(), l)

	cmd := &cobra.Command{
		Use:   "get [key]",
		Short: "Get configuration",
		Args:  cobra.ExactArgs(1),
		// TODO: add argument completion
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Key = args[0]
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
