package use

import (
	"errors"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/profile/ui/cli/shared"
	"github.com/spf13/cobra"
)

// CreateUseCmd constructs and returns the cobra.Command for the profile switch operation.
func CreateUseCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:               "use [name]",
		Short:             "Switch the active configuration profile",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: shared.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]

			if exists, err := container.Profile().Exists(cmd.Context(), opts.Name); err != nil {
				return err
			} else if !exists {
				return coreerrors.NewNotFound(
					errors.New("profile not found"),
					fmt.Sprintf("profile '%s' does not exist", opts.Name),
					"Create a new profile using the 'profile create' command",
				)
			}

			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.OutOrStderr()

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(container.Profile(), container.Logger()).Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
