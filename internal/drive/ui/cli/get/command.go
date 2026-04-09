package get

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/drive/ui/cli/shared"
	"github.com/spf13/cobra"
)

// CreateGetCmd constructs and returns the cobra.Command for the 'drive get' operation.
func CreateGetCmd(container di.Container) *cobra.Command {
	var opts Options

	return &cobra.Command{
		Use:               "get [drive-ref]",
		Short:             "Display details for a specific drive",
		Long:              "Provide a drive ID, name, or alias. If omitted, the currently active drive is shown.",
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: shared.ProviderPathCompletion(container, true),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.DriveRef = args[0]
			opts.Stdout = cmd.OutOrStdout()
			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(container.Drive(), container.Alias(), container.Logger()).Handle(cmd.Context(), opts)
		},
	}
}
