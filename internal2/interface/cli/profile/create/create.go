package create

import (
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const (
	commandName = "create"
	loggerID    = "cli"

	setCurrentFlagLong  = "set-current"
	setCurrentFlagUsage = "Sets the new profile as current"

	forceFlagLong  = "force"
	forceFlagShort = "f"
	forceFlagUsage = "Overwrite an existing profile if it already exists"
)

func CreateCreateCmd(container di.Container) *cobra.Command {
	var (
		setCurrent bool
		force      bool
	)

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new profile",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := util.EnsureLogger(container, loggerID)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			name := strings.ToLower(strings.TrimSpace(args[0]))
			if name == "" {
				logger.Warn("profile name is empty")
				return util.NewCommandErrorWithNameWithMessage(commandName, "name is empty")
			}

			logger.Info("checking if profile exists", infralogging.String("name", name))

			exists, err := container.Profile().Exists(name)
			if err != nil {
				logger.Error("failed to check profile existence", infralogging.String("error", err.Error()))
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			// Handle existing profile
			if exists {
				if !force {
					logger.Warn("profile already exists", infralogging.String("name", name))
					return util.NewCommandErrorWithNameWithMessage(commandName, "profile already exists")
				}

				logger.Warn("profile exists; force enabled, deleting existing profile",
					infralogging.String("name", name),
				)

				if err := container.Profile().Delete(name); err != nil {
					logger.Error("failed to delete existing profile", infralogging.String("error", err.Error()))
					return util.NewCommandErrorWithNameWithError(commandName, err)
				}
			}

			logger.Info("creating profile", infralogging.String("name", name))

			p, err := container.Profile().Create(name)
			if err != nil {
				logger.Error("failed to create profile", infralogging.String("error", err.Error()))
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			if setCurrent {
				logger.Info("setting new profile as current", infralogging.String("name", name))

				if err := container.State().SetCurrentProfile(p.Name); err != nil {
					logger.Error("failed to set current profile", infralogging.String("error", err.Error()))
					return util.NewCommandErrorWithNameWithError(commandName, err)
				}
			}

			logger.Info("profile created successfully",
				infralogging.String("name", p.Name),
				infralogging.String("path", p.Path),
			)

			cmd.Printf("Created profile %q at %s\n", p.Name, p.Path)
			return nil
		},
	}

	cmd.Flags().BoolVar(&setCurrent, setCurrentFlagLong, false, setCurrentFlagUsage)
	cmd.Flags().BoolVarP(&force, forceFlagLong, forceFlagShort, false, forceFlagUsage)

	return cmd
}
