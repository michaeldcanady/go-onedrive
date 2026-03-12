package cat

import (
	"github.com/michaeldcanady/go-onedrive/internal2/di"
	"github.com/spf13/cobra"
)

// CreateCatCmd constructs and returns the cobra.Command for the drive cat operation.
func CreateCatCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "cat path",
		Short: "Display file contents",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Path = args[0]
			opts.Stdout = cmd.OutOrStdout()

			if err := opts.Validate(); err != nil {
				return err
			}

			l, _ := container.Logger().CreateLogger("drive-cat")
			handler := NewHandler(container.FS(), l)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
