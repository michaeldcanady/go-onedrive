package list

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateListCmd constructs and returns the cobra.Command for the drive list operation.
func CreateListCmd(container di.Container) *cobra.Command {
	var opts Options

	return &cobra.Command{
		Use:   "list",
		Short: "List all available OneDrive drives",
		Long: `Retrieve all OneDrive drives associated with your account,
marking the currently active drive and showing any defined aliases.`,
		Args: cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()
			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(container.Drive(), container.Alias(), container.Logger()).Handle(cmd.Context(), opts)
		},
	}
}
