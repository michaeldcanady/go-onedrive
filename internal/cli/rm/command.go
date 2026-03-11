package rm

import (
	"fmt"

	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "rm"
)

func CreateCmd(c didomain.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [path]", commandName),
		Short: "Remove a file or directory from your OneDrive filesystem.",
		Long: `
Remove a file or directory at the specified OneDrive path.

This command behaves similarly to the Unix 'rm' utility, but operates
against your OneDrive domainaccount.

Key behaviors:

  • Removes an item at the given OneDrive path.
  • Fails if the path does not exist.
  • Supports recursive removal of directories with the '--recursive' ('-r') flag.
  • Requires authentication via 'odc auth login'.

Path semantics:
  • All paths refer to locations in your OneDrive filesystem.
  • Absolute paths begin with '/' (recommended).
  • Relative paths are resolved against your OneDrive root.

Authentication:
You must be logged in (via 'odc auth login') before using this command.
`,
		Example: `
  # Remove a file in the root of OneDrive
  odc rm /notes.txt

  # Remove a directory recursively
  odc rm -r /OldProject
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

			rmCmd := NewCmd(c)
			return rmCmd.Run(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "Remove directories and their contents recursively")

	return cmd
}
