package remove

import (
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/spf13/cobra"
)

const (
	commandName = "remove"
	loggerID    = "cli"
)

func CreateRemoveCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "remove <alias>",
		Short: "Remove a drive alias",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Alias = args[0]
			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			c := NewRemoveCmd(container)
			return c.Run(cmd.Context(), opts)
		},
	}

	return cmd
}
