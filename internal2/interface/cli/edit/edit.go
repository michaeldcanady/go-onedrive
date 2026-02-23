package edit

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "edit"
)

// CreateEditCmd constructs and returns the cobra.Command for the edit operation.
// It initializes flags and sets up the execution logic using EditCmd.
func CreateEditCmd(c di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [path]", commandName),
		Short: "Edit a OneDrive file in your local editor",
		Long: `
Edit a OneDrive file in your local editor.

This command downloads the specified file to a temporary local location,
launches your default text editor, and waits for you to save and close it.
If changes are detected (via SHA-256 hashing), the updated file is
automatically uploaded back to OneDrive.

Editor Detection Order:
  1. VISUAL environment variable
  2. EDITOR environment variable
  3. OS Default (notepad on Windows, open -W -t on macOS, xdg-open on Linux)
  4. Common fallbacks (vim, vi, nano)

Authentication:
You must be logged in (via 'onedrive auth login') before using this command.
`,
		Args: cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Path = args[0]
			if err := opts.Validate(); err != nil {
				return err
			}

			opts.Stdin = cmd.InOrStdin()
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			editCmd := NewEditCmd(c)
			return editCmd.Run(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Overwrite the file even if it changed in the cloud")

	return cmd
}
