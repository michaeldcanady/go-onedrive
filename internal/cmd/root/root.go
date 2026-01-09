package root

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/cmd/auth"
	"github.com/michaeldcanady/go-onedrive/internal/cmd/ls"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfigFile     = "./config.yaml"
	loggingLevelFlagLong  = "level"
	loggingLevelFlagShort = "l"
	loggingLevelUsage     = "set the logging level (e.g., debug, info, warn, error)"
)

// CreateRootCmd constructs the root command for the CLI application.
func CreateRootCmd() (*cobra.Command, error) {
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
		return nil, fmt.Errorf("failed to get config flag: %w", err)
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

	rootCmd.PersistentFlags().BoolP(loggingLevelFlagLong, loggingLevelFlagShort, false, loggingLevelUsage)
	if err := viper.BindPFlag("logging.level", rootCmd.PersistentFlags().Lookup(loggingLevelFlagLong)); err != nil {
		return nil, fmt.Errorf("failed to bind logging level flag: %w", err)
	}

	var cfg config.ConfigImpl
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	ctx := context.Background()
	container, err := di.NewContainer(ctx, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize container: %w", err)
	}

	rootCmd.AddCommand(
		ls.CreateLSCmd(container.DriveService, container.Logger),
		auth.CreateAuthCmd(container),
	)

	return rootCmd, nil
}
