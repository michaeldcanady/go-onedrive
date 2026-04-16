package set

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/drive/alias/ui/cli/shared"
	"github.com/spf13/cobra"
)

// CreateSetCmd constructs and returns the cobra.Command for the 'drive alias set' operation.
func CreateSetCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("drive-alias-set")
	handler := NewCommand(container.Alias(), l)

	cmd := &cobra.Command{
		Use:               "set <name> <id>",
		Short:             "Set a drive alias",
		Long:              "Assign a human-readable name to a OneDrive drive ID for easier reference.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: shared.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Alias = args[0]
			opts.DriveID = args[1]
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
