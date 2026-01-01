/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/michaeldcanady/go-onedrive/internal/app"
	"github.com/michaeldcanady/go-onedrive/internal/cache/fsstore"
	jsoncodec "github.com/michaeldcanady/go-onedrive/internal/cache/json_codex"
	"github.com/michaeldcanady/go-onedrive/internal/cmd/ls"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	configFile = "./config.yaml"
)

var (
	graphClient *msgraphsdkgo.GraphServiceClient
	logger      logging.Logger

	ErrMissingAuthConfig    = errors.New("missing 'auth' config section")
	ErrMissingLoggingConfig = errors.New("missing 'logging' config section")
	ErrUnmarshalAuthConfig  = errors.New("unable to unmarshal auth config")
	ErrUnmarshalLogConfig   = errors.New("unable to unmarshal logging config")
)

// rootCmd is the base command for the CLI.
var rootCmd = &cobra.Command{
	Use:   "go-onedrive",
	Short: "A OneDrive CLI client",
	Long: `A command-line interface for interacting with Microsoft OneDrive.

This tool supports authentication, file operations, and integration with
Microsoft Graph APIs.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return writeConfig()
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.new.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	if err := readConfig(); err != nil {
	}
	if err := initializeLogger(); err != nil {
	}
	if err := initializeProfileService(); err != nil {
	}
	if err := initializeCredentialService(); err != nil {
	}
	if err := initializeGraphService(context.Background()); err != nil {
	}

	driveSvc := app.NewDriveService(graphClientService)

	lsCmd := ls.CreateLSCmd(driveSvc)

	rootCmd.AddCommand(lsCmd)
}

func readConfig() error {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		// Missing config file is allowed
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}
		return err
	}
	return nil
}

func writeConfig() error {
	return viper.WriteConfig()
}

func initializeGraphService(ctx context.Context) error {
	sub := viper.Sub("auth")
	if sub == nil {
		return ErrMissingAuthConfig
	}

	var authCfg config.AuthenticationConfigImpl
	if err := sub.Unmarshal(&authCfg); err != nil {
		return errors.Join(ErrUnmarshalAuthConfig, err)
	}

	graphClientService = app.NewGraphClientService(credentialService)

	client, err := graphClientService.Client(ctx)
	if err != nil {
		return fmt.Errorf("failed to create graph client: %w", err)
	}

	graphClient = client
	return nil
}

func initializeLogger() error {
	sub := viper.Sub("logging")
	if sub == nil {
		// TODO: apply default logger config
		return ErrMissingLoggingConfig
	}

	var logCfg config.LoggingConfigImpl
	if err := sub.Unmarshal(&logCfg); err != nil {
		return errors.Join(ErrUnmarshalLogConfig, err)
	}

	cfg := zap.NewProductionConfig()

	switch logCfg.GetLevel() {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return fmt.Errorf("unknown logging level: %s", logCfg.GetLevel())
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("failed to build zap logger: %w", err)
	}

	logger = logging.NewZapLoggerAdapter(zapLogger)
	return nil
}

func initializeProfileService() error {
	store := fsstore.New(".")
	codec := jsoncodec.New()

	profileService = app.NewProfileService(store, codec)
	return nil
}

func initializeCredentialService() error {
	if profileService == nil {
		return errors.New("profile service is nil")
	}

	credentialService = app.NewCredentialService(profileService, logger)
	return nil
}
