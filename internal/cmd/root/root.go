package root

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/cmd/auth"
	"github.com/michaeldcanady/go-onedrive/internal/cmd/ls"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"github.com/spf13/cobra"
)

const (
	profileNameFlagLong     = "profile"
	profileNameFlagDefault  = "default"
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

func CreateRootCmd() (*cobra.Command, error) {
	var (
		level   string
		config  string
		profile string
	)

	ctx := context.Background()

	// Create lightweight container (no heavy services yet)
	container, err := di.NewContainer1(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize container: %w", err)
	}

	// CLI logger (safe to create early)
	cliLogger, err := container.LoggerService.CreateLogger("cli")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cli logger: %w", err)
	}

	rootCmd := &cobra.Command{
		Use:           container.EnvironmentService.Name(),
		Short:         rootShortDescription,
		Long:          rootLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			container.LoggerService.SetAllLevel(level)
			cliLogger.Info("updated all logger level", logging.String("level", level))
			cliLogger.Info("updated config path", logging.String("path", config))
			container.Options = di.RuntimeOptions{
				LogLevel:    level,
				ConfigPath:  config,
				ProfileName: profile,
			}

			return nil
		},
	}

	//
	// ────────────────────────────────────────────────
	// Global Flags
	// ────────────────────────────────────────────────
	//
	rootCmd.PersistentFlags().StringVar(&config, configFileFlagLong, configFileFlagDefault, configFileFlagUsage)
	rootCmd.PersistentFlags().StringVar(&level, loggingLevelFlagLong, loggingLevelFlagDefault, loggingLevelFlagUsage)
	rootCmd.PersistentFlags().StringVar(&profile, profileNameFlagLong, profileNameFlagDefault, profileNameFlagUsage)
	//
	// ────────────────────────────────────────────────
	// Subcommands (DI container passed directly)
	// ────────────────────────────────────────────────
	//
	rootCmd.AddCommand(
		ls.CreateLSCmd(container, cliLogger),
		auth.CreateAuthCmd(container, cliLogger),
	)

	return rootCmd, nil
}
