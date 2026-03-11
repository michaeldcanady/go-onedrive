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
		Long: `You can create a new profile to store your OneDrive configuration and
authentication details. This lets you switch between different accounts
seamlessly.`,
		Example: `  # Create a new profile named 'personal'
  odc profile create personal

  # Create a new profile and set it as the current active profile
  odc profile create work --set-current`,
		Args: cobra.ExactArgs(1),

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
