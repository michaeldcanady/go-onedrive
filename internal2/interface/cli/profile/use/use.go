package use

import (
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const (
	commandName = "use"
)

func CreateUseCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "use <name>",
		Short: "Set current profile",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			name := strings.TrimSpace(args[0])
			if name == "" {
				return util.NewCommandErrorWithNameWithMessage(commandName, "name is empty")
			}

			name = strings.ToLower(name)

			// Validate profile exists
			p, err := container.Profile().Get(cmd.Context(), name)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			// Persist as current profile
			if err := container.State().SetCurrentProfile(p.Name); err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			cmd.Printf("Active profile set to %q\n", p.Name)
			return nil
		},
	}
}
