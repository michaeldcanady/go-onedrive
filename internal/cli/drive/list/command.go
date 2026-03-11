package list

import (
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	commandName = "list"
	loggerID    = "cli"
)

// CreateListCmd constructs and returns the cobra.Command for the list operation.
func CreateListCmd(container didomain.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists available drives",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Stdout: cmd.OutOrStdout(),
				Stderr: cmd.ErrOrStderr(),
			}

			return NewListCmd(container).Run(cmd.Context(), opts)
		},
	}
}
