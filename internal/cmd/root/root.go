package root

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/cmd/auth"
	"github.com/michaeldcanady/go-onedrive/internal/cmd/ls"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfigFile    = "./config.yaml"
	loggingLevelFlagLong = "level"
	loggingLevelUsage    = "set the logging level (e.g., debug, info, warn, error)"
)

// CreateRootCmd constructs the root command for the CLI application.
func CreateRootCmd() (*cobra.Command, error) {
	ctx := context.Background()
	container, err := di.NewContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize container: %w", err)
	}

	cliLogger, err := container.LoggerService.CreateLogger("cli")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cli logger: %w", err)
	}

	rootCmd := &cobra.Command{
		Use:   "odc",
		Short: "A OneDrive CLI client",
		Long: `A command-line interface for interacting with Microsoft OneDrive.

Examples:
  # List files in your OneDrive root
  odc ls

  # Authenticate with OneDrive
  odc auth login

  # Show hidden items in a folder
  odc ls -a Documents
`,
		SilenceUsage:  true,
		SilenceErrors: true,

		// Persist config changes only if something modified Viper state.
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.WriteConfig(); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}
			return nil
		},
	}

	rootCmd.PersistentFlags().String("config", defaultConfigFile, "path to config file")
	configPath, err := rootCmd.PersistentFlags().GetString("config")
	if err != nil {
		cliLogger.Error("failed to get user provided config file", logging.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get config flag: %w", err)
	}
	if strings.TrimSpace(configPath) != "" {
		cliLogger.Debug("user provided config", logging.String("path", configPath))
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Missing config file is fine; other errors are not.
	if err := viper.ReadInConfig(); err != nil {
		// Only warn if the file is missing; fail on parse errors.
		if _, notFound := err.(viper.ConfigFileNotFoundError); !notFound {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
		// TODO: Log that config file was not found, using defaults.
	}

	rootCmd.PersistentFlags().Bool(loggingLevelFlagLong, false, loggingLevelUsage)
	if err := viper.BindPFlag("logging.level", rootCmd.PersistentFlags().Lookup(loggingLevelFlagLong)); err != nil {
		return nil, fmt.Errorf("failed to bind logging level flag: %w", err)
	}

	rootCmd.AddCommand(
		ls.CreateLSCmd(container.DriveService, cliLogger),
		auth.CreateAuthCmd(container, cliLogger),
	)

	return rootCmd, nil
}
