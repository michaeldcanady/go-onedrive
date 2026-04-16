package list

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateListCmd constructs and returns the cobra.Command for the 'drive alias list' operation.
func CreateListCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("drive-alias-list")
	handler := NewCommand(container.Alias(), l)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all defined drive aliases",
		Long:  `Retrieve and display a list of all custom aliases defined for OneDrive drives.`,
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
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

