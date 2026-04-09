package current

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateCurrentCmd constructs and returns the cobra.Command for showing the current profile.
func CreateCurrentCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "current",
		Short: "Display the name of the currently active profile",
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
