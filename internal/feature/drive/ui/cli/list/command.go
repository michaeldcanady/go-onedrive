package list

import (
	"github.com/michaeldcanady/go-onedrive/internal/feature/di"
	"github.com/spf13/cobra"
)

// CreateListCmd constructs and returns the cobra.Command for the drive list operation.
func CreateListCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available OneDrive drives",
		Long: `Retrieve all OneDrive drives associated with your account,
marking the currently active drive and showing any defined aliases.`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := NewOptions()
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			log, err := container.Logger().CreateLogger("drive-list")
			if err != nil {
				return err
			}

			handler := NewHandler(container.Drive(), container.State(), log)
			return handler.Handle(cmd.Context(), opts)
		},
	}
}
