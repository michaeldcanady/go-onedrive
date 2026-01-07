package root

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/cmd/auth"
	"github.com/michaeldcanady/go-onedrive/internal/cmd/ls"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// TODO: make an option
	configFile = "./config.yaml"
)

func CreateRootCmd() (*cobra.Command, error) {
	// 1. Load config
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	_ = viper.ReadInConfig() // ignore missing file

	var cfg config.ConfigImpl
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// 2. Build DI container
	ctx := context.Background()
	container, err := di.NewContainer(ctx, &cfg)
	if err != nil {
		return nil, err
	}

	// 3. Build root command
	rootCmd := &cobra.Command{
		Use:   "go-onedrive",
		Short: "A OneDrive CLI client",
		Long:  "A command-line interface for interacting with Microsoft OneDrive.",
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return viper.WriteConfig()
		},
	}

	// 4. Add subcommands with dependencies
	rootCmd.AddCommand(
		ls.CreateLSCmd(container.DriveService, container.Logger),
		auth.CreateAuthCmd(container),
	)

	return rootCmd, nil
}
