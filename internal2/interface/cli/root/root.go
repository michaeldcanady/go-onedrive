package root

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	infraprofile "github.com/michaeldcanady/go-onedrive/internal2/infra/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/auth"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/ls"
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
	rootCmd.PersistentFlags().StringVar(&profile, profileNameFlagLong, infraprofile.DefaultProfileName, profileNameFlagUsage)

	rootCmd.AddCommand(
		ls.CreateLSCmd(container),
		auth.CreateAuthCmd(container),
		profilecmd.CreateProfileCmd(container),
	)

	return rootCmd, nil
}
