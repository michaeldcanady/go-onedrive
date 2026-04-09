package use

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/drive/ui/cli/shared"
	"github.com/spf13/cobra"
)

// CreateUseCmd constructs and returns the cobra.Command for the 'drive use' operation.
func CreateUseCmd(container di.Container) *cobra.Command {
	var opts Options

	return &cobra.Command{
		Use:               "use <drive-ref>",
		Short:             "Set the active OneDrive drive",
		Long:              "Specify a drive ID, name, or alias to be used for subsequent commands.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: shared.ProviderPathCompletion(container, true),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// TODO: should be able to use to use alias to set the active drive, but this will require a lookup to resolve the alias to a drive ID before validation
			opts.DriveRef = args[0]
			opts.Stdout = cmd.OutOrStdout()
			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(container.Drive(), container.Alias(), container.Logger()).Handle(cmd.Context(), opts)
		},
	}
}
