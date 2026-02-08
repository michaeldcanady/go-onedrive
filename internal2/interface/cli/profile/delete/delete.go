package delete

import (
	"context"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
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
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
				cmd.SetContext(ctx)
			}

			logger, err := util.EnsureLogger(ctx, container, loggerID)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("command started",
				logging.String("command", commandName),

				logging.Bool("force", force),
				logging.Strings("args", args),
			)

			name := strings.ToLower(strings.TrimSpace(args[0]))
			if name == "" {
				logger.Warn("profile name is empty",
					logging.String("event", "validate_input"),
				)
				return util.NewCommandErrorWithNameWithMessage(commandName, "name is empty")
			}

			if name == profile.DefaultProfileName {
				logger.Warn("attempted to delete default profile",
					logging.String("event", "validate_input"),
					logging.String("profile", name),
				)
				return util.NewCommandErrorWithNameWithMessage(
					commandName,
					"cannot delete the default profile",
				)
			}

			logger.Debug("retrieving current profile",
				logging.String("event", "resolve_current_profile"),
			)

			current, err := container.State().GetCurrentProfile()
			if err != nil {
				logger.Warn("failed to retrieve current profile",
					logging.String("event", "resolve_current_profile"),
					logging.Error(err),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			// If deleting the active profile, confirm unless forced
			if current == name && !force {
				logger.Info("deleting active profile; confirmation required",
					logging.String("event", "confirm_delete_active"),
					logging.String("profile", name),
				)

				prompt := promptui.Prompt{
					Label:     fmt.Sprintf("You are deleting the active profile %q. Continue", name),
					IsConfirm: true,
				}

				_, err := prompt.Run()
				if err != nil {
					logger.Info("delete aborted by user",
						logging.String("event", "confirm_delete_active"),
						logging.String("profile", name),
					)
					cmd.Println("Aborted.")
					return nil
				}

				logger.Info("switching to default profile before deletion",
					logging.String("event", "switch_profile"),
					logging.String("from_profile", name),
					logging.String("to_profile", profile.DefaultProfileName),
				)

				if err := container.State().SetCurrentProfile(profile.DefaultProfileName); err != nil {
					logger.Warn("failed to switch to default profile",
						logging.String("event", "switch_profile"),
						logging.Error(err),
					)
					return util.NewCommandErrorWithNameWithError(
						commandName,
						fmt.Errorf("failed to switch to default profile: %w", err),
					)
				}
			}

			logger.Info("deleting profile directory",
				logging.String("event", "delete_profile"),
				logging.String("profile", name),
			)

			if err := container.Profile().Delete(name); err != nil {
				logger.Warn("failed to delete profile",
					logging.String("event", "delete_profile"),
					logging.String("profile", name),
					logging.Error(err),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("profile deleted successfully",
				logging.String("event", "delete_success"),
				logging.String("profile", name),
			)

			cmd.Printf("Deleted profile %q\n", name)

			logger.Info("command completed",
				logging.String("command", commandName),
				logging.String("profile", name),
				logging.Bool("force", force),
			)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force deletion without confirmation")

	return cmd
}
