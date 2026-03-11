package create

import (
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	commandName = "create"
	loggerID    = "cli"

	setCurrentFlagLong  = "set-current"
	setCurrentFlagUsage = "Sets the new profile as current"

	forceFlagLong  = "force"
	forceFlagShort = "f"
	forceFlagUsage = "Overwrite an existing profile if it already exists"
)

// CreateCreateCmd constructs and returns the cobra.Command for the create operation.
func CreateCreateCmd(container didomain.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new profile",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			return NewCreateCmd(container).Run(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.SetCurrent, setCurrentFlagLong, false, setCurrentFlagUsage)
	cmd.Flags().BoolVarP(&opts.Force, forceFlagLong, forceFlagShort, false, forceFlagUsage)

	return cmd
}
