package list

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateListCmd constructs and returns the cobra.Command for listing profiles.
func CreateListCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configuration profiles",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(container.Profile(), container.Logger()).Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
