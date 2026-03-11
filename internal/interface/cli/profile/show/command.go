package show

import (
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	commandName = "show"
	loggerID    = "cli"
)

// CreateShowCmd constructs and returns the cobra.Command for the show operation.
func CreateShowCmd(container didomain.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "show <NAME>",
		Short: "Show details for a profile",
		Long: `You can view the detailed configuration and authentication information for a
specific profile. Profiles let you manage multiple OneDrive accounts or
configurations easily.`,
		Example: `  # Show details for the 'personal' profile
  odc profile show personal

  # Show details for the 'work' profile
  odc profile show work`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Name:   args[0],
				Stdout: cmd.OutOrStdout(),
				Stderr: cmd.ErrOrStderr(),
			}

			return NewShowCmd(container).Run(cmd.Context(), opts)
		},
	}
}
