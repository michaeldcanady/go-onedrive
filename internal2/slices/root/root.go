// Package root provides the entry point for the odc2 command-line interface.
package root

import (
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/core/state"
	"github.com/michaeldcanady/go-onedrive/internal2/di"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/auth"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/profile"
	"github.com/spf13/cobra"
)

const (
	rootShortDescription = "A OneDrive CLI client (v2)"
	rootLongDescription  = `A command-line interface for interacting with Microsoft OneDrive, implemented with Vertical Slice Architecture.`
)

// CreateRootCmd constructs and returns the cobra.Command for the root application.
func CreateRootCmd(container di.Container) (*cobra.Command, error) {
	var (
		level   string
		config  string
		profileFlag string
	)

	rootCmd := &cobra.Command{
		Use:           "odc2",
		Short:         rootShortDescription,
		Long:          rootLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cliLogger, err := container.Logger().CreateLogger("cli")
			if err != nil {
				return fmt.Errorf("failed to initialize cli logger: %w", err)
			}

			// Ensure a default profile exists
			if exists, _ := container.Profile().Exists(cmd.Context(), "default"); !exists {
				_, _ = container.Profile().Create(cmd.Context(), "default")
			}

			if strings.TrimSpace(profileFlag) != "" {
				if err := container.State().Set(state.KeyProfile, profileFlag, state.ScopeSession); err != nil {
					return fmt.Errorf("failed to set session profile: %w", err)
				}
			} else {
				// Use default profile
				_ = container.State().Set(state.KeyProfile, "default", state.ScopeSession)
			}

			profileName, err := container.State().Get(state.KeyProfile)
			if err != nil {
				return fmt.Errorf("failed to get current profile: %w", err)
			}

			if err := container.Config().AddPath(profileName, config); err != nil {
				return fmt.Errorf("failed to load config file %s: %w", config, err)
			}

			container.Logger().SetAllLevel(level)
			cliLogger.Debug("updated all logger level", logger.String("level", level))
			cliLogger.Debug("updated config path", logger.String("path", config))

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&config, "config", "./config.yaml", "path to config file")
	rootCmd.PersistentFlags().StringVar(&level, "level", "info", "set the logging level (e.g., debug, info, warn, error)")
	rootCmd.PersistentFlags().StringVar(&profileFlag, "profile", "", "name of profile")

	rootCmd.AddCommand(
		auth.CreateAuthCmd(container),
		profile.CreateProfileCmd(container),
	)

	return rootCmd, nil
}
