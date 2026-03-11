package mv

import (
	"fmt"

	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "mv"
)

func CreateCmd(c didomain.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [source] [destination]", commandName),
		Short: "Move or rename a file/directory in your OneDrive filesystem.",
		Long: `
Move or rename an item at the specified OneDrive source path to the
destination path.

This command behaves similarly to the Unix 'mv' utility, but operates
against your OneDrive domainaccount.

Key behaviors:

  • Moves an item from source to destination.
  • Renames an item if the destination is in the same directory but has a different name.
  • Fails if the source does not exist.
  • Fails if the destination parent directory does not exist.
  • Requires authentication via 'odc auth login'.

Path semantics:
  • All paths refer to locations in your OneDrive filesystem.
  • Absolute paths begin with '/' (recommended).
  • Relative paths are resolved against your OneDrive root.

Authentication:
You must be logged in (via 'odc auth login') before using this command.
`,
		Example: `
  # Rename a file
  odc mv /oldname.txt /newname.txt

  # Move a file to a different directory
  odc mv /file.txt /Documents/file.txt

  # Move and rename a file
  odc mv /file.txt /Documents/newfile.txt
`,

		Args: cobra.ExactArgs(2),

		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Source = args[0]
			opts.Destination = args[1]
			return opts.Validate()
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Stdin = cmd.InOrStdin()
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			mvCmd := NewCmd(c)
			return mvCmd.Run(cmd.Context(), opts)
		},
	}

	return cmd
}
