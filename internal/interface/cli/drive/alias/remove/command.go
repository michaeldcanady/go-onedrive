package remove

import (
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	commandName = "remove"
	loggerID    = "cli"
)

func CreateRemoveCmd(container didomain.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "remove <ALIAS>",
		Short: "Remove a drive alias",
		Long: `You can remove an existing drive alias. This doesn't delete the drive itself;
it only removes the short name you've created for it.`,
		Example: `  # Remove the alias 'personal'
  odc drive alias remove personal`,
		Args: cobra.ExactArgs(1),
		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Alias = args[0]
			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			c := NewRemoveCmd(container)
			return c.Run(cmd.Context(), opts)
		},
	}

	return cmd
}
