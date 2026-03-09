package touch

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "touch"
)

func CreateCmd(c di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [path]", commandName),
		Short: "Create an empty file in your OneDrive filesystem.",
		Long: `
Create a new empty file at the specified OneDrive path.

This command behaves similarly to the Unix 'touch' utility (creation only),
but operates against your OneDrive account.

Key behaviors:

  • Creates a single empty file at the given OneDrive path.
  • Fails if the parent directory does not exist.
  • Requires authentication via 'odc auth login'.

Path semantics:
  • All paths refer to locations in your OneDrive filesystem.
  • Absolute paths begin with '/' (recommended).
  • Relative paths are resolved against your OneDrive root.

Authentication:
You must be logged in (via 'odc auth login') before using this command.
`,
		Example: `
  # Create a file in the root of OneDrive
  odc touch /newfile.txt

  # Create a file inside an existing directory
  odc touch /Documents/notes.md
`,

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
