package root

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/cmd/auth"
	"github.com/michaeldcanady/go-onedrive/internal/cmd/ls"
	"github.com/michaeldcanady/go-onedrive/internal/di"
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
	var container *di.Container

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
			container.ConfigurationService.WriteConfiguration(context.Background())
			return nil
		},
	}

	rootCmd.PersistentFlags().String("config", defaultConfigFile, "path to config file")
	configPath, err := rootCmd.PersistentFlags().GetString("config")
	if err != nil {
		return nil, fmt.Errorf("failed to get config flag: %w", err)
	}

	rootCmd.PersistentFlags().Bool(loggingLevelFlagLong, false, loggingLevelUsage)
	// TODO: how to handle with config service?
	if err := viper.BindPFlag("logging.level", rootCmd.PersistentFlags().Lookup(loggingLevelFlagLong)); err != nil {
		return nil, fmt.Errorf("failed to bind logging level flag: %w", err)
	}

	ctx := context.Background()
	container, err = di.NewContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize container: %w", err)
	}

	container.ConfigurationService.SetConfigFile(ctx, configPath)
	rootCmd.AddCommand(
		ls.CreateLSCmd(container.DriveService, container.Logger),
		auth.CreateAuthCmd(container),
	)

	return rootCmd, nil
}
