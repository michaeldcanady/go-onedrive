package cp

import (
	"fmt"

	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "cp"
)

// CreateCpCmd constructs and returns the cobra.Command for the cp operation.
// It initializes flags and sets up the execution logic using CpCmd.
func CreateCpCmd(c didomain.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s <source> <dest>", commandName),
		Short: "Copy items from source to destination",
		Long: `You can copy items from a source path to a destination path using this command.
It supports local-to-remote, remote-to-local, and remote-to-remote copying.
You can prefix paths with 'local:' or 'onedrive:' (the default).`,
		Example: `  # Copy a local file to OneDrive
  odc cp local:report.docx onedrive:/Documents/report.docx

  # Copy a OneDrive file to your local machine
  odc cp onedrive:/Photos/vacation.jpg local:vacation.jpg

  # Copy a file within OneDrive
  odc cp /Projects/budget.xlsx /Archive/budget_v1.xlsx`,

		Args: cobra.ExactArgs(2),

		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Source = args[0]
			opts.Dest = args[1]
			return opts.Validate()
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Stdin = cmd.InOrStdin()
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			runCmd := NewCpCmd(c)
			return runCmd.Run(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Overwrite, "overwrite", "f", false, "Overwrite destination if it exists")
	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "Recursively copy all files within a directory")

	return cmd
}
