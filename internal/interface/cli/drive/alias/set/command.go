package set

import (
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	commandName = "set"
	loggerID    = "cli"
)

func CreateSetCmd(container didomain.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "set <ALIAS> <DRIVE_ID>",
		Short: "Set a drive alias",
		Long: `You can create a convenient alias for a long OneDrive drive ID. This makes it
much easier to refer to your drives in other commands. If the alias already
exists, it's updated to point to the new drive ID.`,
		Example: `  # Set an alias for your personal drive
  odc drive alias set personal b!1234567890abcdef

  # Set an alias for a shared project drive
  odc drive alias set project-x b!0987654321fedcba`,
		Args: cobra.ExactArgs(2),
		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Alias = args[0]
			opts.DriveID = args[1]
			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			c := NewSetCmd(container)
			return c.Run(cmd.Context(), opts)
		},
	}

	return cmd
}
