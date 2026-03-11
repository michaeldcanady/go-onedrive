package list

import (
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	commandName = "list"
	loggerID    = "cli"
)

func CreateListCmd(container didomain.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List drive aliases",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(_ *cobra.Command, args []string) error {
			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			c := NewListCmd(container)
			return c.Run(cmd.Context(), opts)
		},
	}

	return cmd
}
