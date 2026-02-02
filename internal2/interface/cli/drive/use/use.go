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
		Use:   "use <id>",
		Short: "Set current drive",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: add logger

			id := strings.TrimSpace(args[0])
			if id == "" {
				return util.NewCommandErrorWithNameWithMessage(commandName, "id is empty")
			}

			id = strings.ToLower(id)

			// TODO: update drive state
			_, err := container.Drive().ResolveDrive(cmd.Context(), id)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			return nil
		},
	}
}
