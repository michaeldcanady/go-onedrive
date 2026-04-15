package list

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateListCmd constructs and returns the cobra.Command for the drive list operation.
func CreateListCmd(container di.Container) *cobra.Command {
	var opts Options

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

			return handler.Validate(cmd.Context(), &opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return handler.Execute(cmd.Context(), opts)
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return handler.Finalize(cmd.Context(), opts)
		},
	}

	return cmd
}
