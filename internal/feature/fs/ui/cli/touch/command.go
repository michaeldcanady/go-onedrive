package touch

import (
	"github.com/michaeldcanady/go-onedrive/internal/feature/di"
	"github.com/spf13/cobra"
)

// CreateTouchCmd constructs and returns the cobra.Command for the drive touch operation.
func CreateTouchCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "touch <path>",
		Short: "Create an empty file or update its timestamp",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Path = args[0]
			opts.Stdout = cmd.OutOrStdout()

			if err := opts.Validate(); err != nil {
				return err
			}

			l, _ := container.Logger().CreateLogger("drive-touch")
			handler := NewHandler(container.FS(), l)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
