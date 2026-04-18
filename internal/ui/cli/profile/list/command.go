package list

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateListCmd constructs and returns the cobra.Command for the profile list operation.
func CreateListCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("profile-list")
	handler := NewHandler(container.Profile(), l)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all profiles",
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
			return handler.Execute(c)
		},

		PostRunE: func(cmd *cobra.Command, args []string) error {
			return handler.Finalize(c)
		},
	}

	return cmd
}
