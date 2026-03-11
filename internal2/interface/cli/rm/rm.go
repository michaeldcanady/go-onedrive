package rm

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "rm"
)

func CreateCmd(c di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [path]", commandName),
		Short: "Remove a file or directory from your OneDrive filesystem.",
		Long:  "Remove a file or directory from your OneDrive filesystem. By default, items are moved to the recycle bin.",
		Example: `  # Remove a file (moves to recycle bin)
  odc rm /path/to/file.txt

  # Permanently remove a file
  odc rm --permanent /path/to/file.txt`,

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

	// TODO: need to open bug with MS Graph SDk https://github.com/microsoftgraph/msgraph-sdk-go/issues/980
	//cmd.Flags().BoolVar(&opts.Permanent, "permanent", false, "Permanently delete the item instead of moving it to the recycle bin")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Force permanent deletion without confirmation prompt")

	return cmd
}
