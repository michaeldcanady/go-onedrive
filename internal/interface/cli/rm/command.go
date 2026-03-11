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
		Use:   fmt.Sprintf("%s <PATH>", commandName),
		Short: "Remove a file or directory from your OneDrive filesystem.",
		Long: `You can remove files or directories from your OneDrive. This command works
similarly to the Unix 'rm' utility. You can use the recursive flag to delete
entire directories and their contents.`,
		Example: `  # Delete a single file from OneDrive
  odc rm /Old/temp.txt

  # Delete a folder and everything inside it
  odc rm -r /Backups/2023`,

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
