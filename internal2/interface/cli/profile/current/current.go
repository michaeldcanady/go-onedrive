package current

import (
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const commandName = "current"

func CreateCurrentCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show the active profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := container.State().GetCurrentProfile()
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			cmd.Printf("%s\n", name)
			return nil
		},
	}
}
