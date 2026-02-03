package delete

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const (
	commandName = "delete"
	loggerID    = "cli"
)

func CreateDeleteCmd(container di.Container) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a profile",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := util.EnsureLogger(container, loggerID)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			name := strings.ToLower(strings.TrimSpace(args[0]))
			if name == "" {
				return util.NewCommandErrorWithNameWithMessage(commandName, "name is empty")
			}

			if name == profile.DefaultProfileName {
				return util.NewCommandErrorWithNameWithMessage(
					commandName,
					"cannot delete the default profile",
				)
			}

			current, err := container.State().GetCurrentProfile()
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			// If deleting the active profile, confirm unless forced
			if current == name && !force {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf("You are deleting the active profile %q. Continue", name),
					IsConfirm: true,
				}

				_, err := prompt.Run()
				if err != nil {
					cmd.Println("Aborted.")
					return nil
				}

				logger.Info("deleting current profile; switching to default")

				if err := container.State().SetCurrentProfile(profile.DefaultProfileName); err != nil {
					return util.NewCommandErrorWithNameWithError(
						commandName,
						fmt.Errorf("failed to switch to default profile: %w", err),
					)
				}
			}

			// Delete the profile directory
			if err := container.Profile().Delete(name); err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			cmd.Printf("Deleted profile %q\n", name)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force deletion without confirmation")

	return cmd
}
