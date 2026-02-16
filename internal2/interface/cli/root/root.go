package root

import (
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	authcmd "github.com/michaeldcanady/go-onedrive/internal2/interface/cli/auth"
	catcmd "github.com/michaeldcanady/go-onedrive/internal2/interface/cli/cat"
	drivecmd "github.com/michaeldcanady/go-onedrive/internal2/interface/cli/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/edit"
	lscmd "github.com/michaeldcanady/go-onedrive/internal2/interface/cli/ls"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/middleware"
	profilecmd "github.com/michaeldcanady/go-onedrive/internal2/interface/cli/profile"
	"github.com/spf13/cobra"
)

const (
	profileNameFlagLong     = "profile"
	profileNameFlagUsage    = "name of profile"
	configFileFlagLong      = "config"
	configFileFlagUsage     = "path to config file"
	configFileFlagDefault   = "./config.yaml"
	loggingLevelFlagLong    = "level"
	loggingLevelFlagUsage   = "set the logging level (e.g., debug, info, warn, error)"
	loggingLevelFlagDefault = "info"
	rootShortDescription    = "A OneDrive CLI client"
	rootLongDescription     = `A command-line interface for interacting with Microsoft OneDrive.

Examples:
  # List files in your OneDrive root
  odc ls

  # Authenticate with OneDrive
  odc auth login

  # Show hidden items in a folder
  odc ls -a Documents
`
)

func CreateRootCmd(container di.Container) (*cobra.Command, error) {
	var (
		level   string
		config  string
		profile string
	)

	rootCmd := &cobra.Command{
		Use:           container.EnvironmentService().Name(),
		Short:         rootShortDescription,
		Long:          rootLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// CLI logger (safe to create early)
			cliLogger, err := container.Logger().CreateLogger("cli")
			if err != nil {
				return fmt.Errorf("failed to initialize cli logger: %w", err)
			}

			if strings.TrimSpace(profile) != "" {
				container.State().SetSessionProfile(profile)
			}

			profileName, err := container.State().GetCurrentProfile()
			if err != nil {
				return fmt.Errorf("failed to get current profile: %w", err)
			}

			if err := container.Config().AddPath(profileName, config); err != nil {
				return fmt.Errorf("failed to load config file %s: %w", config, err)
			}

			container.Logger().SetAllLevel(level)
			cliLogger.Info("updated all logger level", logging.String("level", level))
			cliLogger.Info("updated config path", logging.String("path", config))

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&config, configFileFlagLong, configFileFlagDefault, configFileFlagUsage)
	rootCmd.PersistentFlags().StringVar(&level, loggingLevelFlagLong, loggingLevelFlagDefault, loggingLevelFlagUsage)
	rootCmd.PersistentFlags().StringVar(&profile, profileNameFlagLong, "", profileNameFlagUsage)

	rootCmd.AddCommand(
		lscmd.CreateLSCmd(container),
		authcmd.CreateAuthCmd(container),
		profilecmd.CreateProfileCmd(container),
		drivecmd.CreateDriveCmd(container),
		catcmd.CreateCatCmd(container),
		edit.CreateEditCmd(container),
	)

	middleware.ApplyMiddlewareRecursively(rootCmd, middleware.WithCorrelationID)

	return rootCmd, nil
}
