package set

import (
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/spf13/cobra"
)

const (
	commandName = "set"
	loggerID    = "cli"
)

func CreateSetCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "set <alias> <drive-id>",
		Short: "Set a drive alias",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Alias = args[0]
			opts.DriveID = args[1]
			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			c := NewCmd(container)
			return c.Run(cmd.Context(), opts)
		},
	}

	return cmd
}
