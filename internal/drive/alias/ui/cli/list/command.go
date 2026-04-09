package list

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateListCmd constructs and returns the cobra.Command for the 'drive alias list' operation.
func CreateListCmd(container di.Container) *cobra.Command {
	var opts Options

	return &cobra.Command{
		Use:   "list",
		Short: "List all drive aliases",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()
			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(container.Alias(), container.Logger()).Handle(cmd.Context(), opts)
		},
	}
}
