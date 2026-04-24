package get

import (
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/spf13/cobra"
)

// CreateGetCmd constructs and returns the cobra.Command for the 'drive get' operation.
func CreateGetCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("drive-get")
	handler := NewCommand(container.Drive(), l)

	cmd := &cobra.Command{
		Use:   "get <drive-ref>",
		Short: "Display details for a specific drive",
		Long:  "Retrieve and show the metadata for a OneDrive drive identified by its ID or name.",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.DriveRef = args[0]
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

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

	cmd.Flags().StringVar(&opts.IdentityID, "id", "", "The specific identity (email or alias) to get the personal drive for")

	return cmd
}
