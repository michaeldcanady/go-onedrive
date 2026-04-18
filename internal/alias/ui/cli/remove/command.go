package remove

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/alias/ui/cli/shared"
	"github.com/spf13/cobra"
)

// CreateRemoveCmd constructs and returns the cobra.Command for the 'drive alias remove' operation.
func CreateRemoveCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("drive-alias-remove")
	handler := NewCommand(container.Alias(), l)

	cmd := &cobra.Command{
		Use:               "remove <name>",
		Short:             "Remove a drive alias",
		Long:              "Delete a previously defined drive alias.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: shared.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Alias = args[0]
			opts.Stdout = cmd.OutOrStdout()

			c = &CommandContext{
				Ctx:     cmd.Context(),
				Options: opts,
			}

			return handler.Validate(c)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := handler.Execute(c); err != nil {
				return err
			}
			return handler.Finalize(c)
		},
	}

	return cmd
}
