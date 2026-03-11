package touch

import (
	"fmt"

	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "touch"
)

func CreateCmd(c didomain.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [path]", commandName),
		Short: "Create an empty file in your OneDrive filesystem.",
		Long: `You can create a new empty file at the specified OneDrive path. This command
works like the Unix 'touch' utility, but it operates on your OneDrive
account. It's useful for initializing files or creating placeholders.`,
		Example: `  # Create a new empty file in the OneDrive root
  odc touch /new_project.txt

  # Create a file inside a specific folder
  odc touch /Documents/notes/meeting_notes.md`,

		Args: cobra.ExactArgs(1),

		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Path = args[0]
			return opts.Validate()
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Stdin = cmd.InOrStdin()
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			touchCmd := NewCmd(c)
			return touchCmd.Run(cmd.Context(), opts)
		},
	}

	return cmd
}
