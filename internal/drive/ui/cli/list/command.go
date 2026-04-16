package list

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateListCmd constructs and returns the cobra.Command for the drive list operation.
func CreateListCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("drive-list")
	handler := NewCommand(container.Drive(), container.Alias(), l)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available OneDrive drives",
		Long: `Retrieve all OneDrive drives associated with your account,
marking the currently active drive and showing any defined aliases.`,
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
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

	return cmd
}
