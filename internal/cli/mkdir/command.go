package mkdir

import (
	"fmt"

	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "mkdir"

	parentCommandFlagLong  = "parent"
	ParentCommandFlagShort = "p"
)

func CreateCmd(c didomain.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [path]", commandName),
		Short: "Create a directory in your OneDrive filesystem.",
		Long: `You can create a new directory in your OneDrive. By default, the parent
directory must already exist unless you use the 'parent' flag to create any
missing folders in the path.`,
		Example: `  # Create a new directory named 'Photos' in the root
  odc mkdir /Photos

  # Create a nested directory path, including any missing parents
  odc mkdir -p /Work/Projects/2024/Finances`,

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
