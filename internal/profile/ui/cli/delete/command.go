package delete

import (
	"errors"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/profile/ui/cli/shared"
	"github.com/spf13/cobra"
)

// CreateDeleteCmd constructs and returns the cobra.Command for the profile deletion operation.
func CreateDeleteCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:               "delete [name]",
		Short:             "Delete a configuration profile",
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
					"Profile already does not exist, no need to delete it",
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
