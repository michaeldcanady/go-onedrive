package cp

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "cp"
)

// CreateCpCmd constructs and returns the cobra.Command for the cp operation.
// It initializes flags and sets up the execution logic using CpCmd.
func CreateCpCmd(c di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s <source> <dest>", commandName),
		Short: "Copy items from source to destination",
		Long: `
Copy items from a source path to a destination path.
Supports local-to-remote, remote-to-local, and remote-to-remote copying.
Paths can be prefixed with 'local:' or 'onedrive:' (default).
`,
		Example: `
  # Copy local file to OneDrive
  odc cp local:file.txt onedrive:/folder/file.txt

  # Copy OneDrive folder to local
  odc cp -r onedrive:/folder local:./folder

  # Copy within OneDrive
  odc cp /folder/file1.txt /folder/file2.txt
`,

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
	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "Copy directories recursively")

	return cmd
}
