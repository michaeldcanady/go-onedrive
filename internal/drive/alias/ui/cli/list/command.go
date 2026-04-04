package list

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateListCmd constructs and returns the cobra.Command for the 'drive alias list' operation.
func CreateListCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all drive aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{Stdout: cmd.OutOrStdout()}
			log, _ := container.Logger().CreateLogger("alias-list")
			return NewHandler(container.Alias(), log).Handle(cmd.Context(), opts)
		},
	}
}
