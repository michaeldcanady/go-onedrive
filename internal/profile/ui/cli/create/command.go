package create

import (
	"errors"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/spf13/cobra"
)

// CreateCreateCmd constructs and returns the cobra.Command for the profile creation operation.
func CreateCreateCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new configuration profile",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]

			if exists, err := container.Profile().Exists(cmd.Context(), opts.Name); err != nil {
				return err
			} else if exists {
				return coreerrors.NewConflict(
					errors.New("profile already exists"),
					fmt.Sprintf("profile '%s' already exists", opts.Name),
					"Use the 'profile use' command to switch to an existing profile",
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
