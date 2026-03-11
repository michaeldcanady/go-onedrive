package cat

import (
	"fmt"

	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "cat"
)

// CreateCatCmd constructs and returns the cobra.Command for the cat operation.
func CreateCatCmd(c didomain.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [path]", commandName),
		Short: "Display the contents of a OneDrive file",
		Args:  cobra.ExactArgs(1),

		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Path = args[0]
			return opts.Validate()
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Stdin = cmd.InOrStdin()
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			catCmd := NewCatCmd(c)
			return catCmd.Run(cmd.Context(), opts)
		},
	}

	return cmd
}
