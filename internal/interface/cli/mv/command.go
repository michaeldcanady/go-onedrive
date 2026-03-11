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
		Use:   fmt.Sprintf("%s <SOURCE> <DESTINATION>", commandName),
		Short: "Move or rename a file/directory in your OneDrive filesystem.",
		Long: `You can move or rename files and directories within your OneDrive. This
command works like the Unix 'mv' utility, allowing you to change an item's
location, its name, or both at once.`,
		Example: `  # Rename a file in the same directory
  odc mv /old_report.docx /final_report.docx

  # Move a file to a different folder
  odc mv /drafts/memo.txt /Archive/memo.txt

  # Move and rename a folder
  odc mv /Projects/Current /Projects/Archived/ProjectAlpha`,

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
