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
		Use:     fmt.Sprintf("%s [path]", commandName),
		Short:   "Move or rename a file/directory in your OneDrive filesystem.",
		Long:    "",
		Example: "",

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

	return cmd
}
