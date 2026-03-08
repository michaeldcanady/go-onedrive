package mkdir

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "mkdir"

	parentCommandFlagLong  = "parent"
	ParentCommandFlagShort = "p"
)

func CreateCmd(c di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [path]", commandName),
		Short: "Create a directory in your OneDrive filesystem.",
		Long: `
Create a new directory at the specified OneDrive path.

This command behaves similarly to the Unix 'mkdir' utility, but operates
against your OneDrive account. By default, the parent directory must already
exist. Use '--parent' ('-p') to create any missing parent directories
recursively.

Key behaviors:

  • Creates a single directory at the given OneDrive path.
  • Fails if the parent directory does not exist, unless '--parent' is used.
  • Succeeds without error if the directory already exists and '--parent' is set.
  • Requires authentication via 'odc auth login'.

Path semantics:
  • All paths refer to locations in your OneDrive filesystem.
  • Absolute paths begin with '/' (recommended).
  • Relative paths are resolved against your OneDrive root.

Authentication:
You must be logged in (via 'odc auth login') before using this command.
`,
		Example: `
  # Create a directory in the root of OneDrive
  odc mkdir /Projects

  # Create nested directories, creating parents as needed
  odc mkdir -p /Projects/2025/Reports

  # Create a folder inside an existing directory
  odc mkdir /Documents/Invoices
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

			mkdirCmd := NewCmd(c)
			return mkdirCmd.Run(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Parent, parentCommandFlagLong, ParentCommandFlagShort, false, "Recursively create parent directories")

	return cmd
}
