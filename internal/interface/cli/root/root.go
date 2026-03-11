package root

import (
	"fmt"
	"strings"

	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"

	authcmd "github.com/michaeldcanady/go-onedrive/internal/interface/cli/auth"
	catcmd "github.com/michaeldcanady/go-onedrive/internal/interface/cli/cat"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/cp"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/download"
	drivecmd "github.com/michaeldcanady/go-onedrive/internal/interface/cli/drive"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/edit"
	lscmd "github.com/michaeldcanady/go-onedrive/internal/interface/cli/ls"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/middleware"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/mkdir"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/mv"
	profilecmd "github.com/michaeldcanady/go-onedrive/internal/interface/cli/profile"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/rm"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/touch"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/upload"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
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

func CreateRootCmd(container didomain.Container) (*cobra.Command, error) {
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
				if err := container.State().Set(domainstate.KeyProfile, profile, domainstate.ScopeSession); err != nil {
					return fmt.Errorf("failed to set session profile: %w", err)
				}
			}

			profileName, err := container.State().Get(domainstate.KeyProfile)
			if err != nil {
				return fmt.Errorf("failed to get current profile: %w", err)
			}

			if err := container.Config().AddPath(profileName, config); err != nil {
				return fmt.Errorf("failed to load config file %s: %w", config, err)
			}

			container.Logger().SetAllLevel(level)
			cliLogger.Info("updated all logger level", domainlogger.String("level", level))
			cliLogger.Info("updated config path", domainlogger.String("path", config))

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
		upload.CreateUploadCmd(container),
		download.CreateDownloadCmd(container),
		mkdir.CreateCmd(container),
		touch.CreateCmd(container),
		mv.CreateCmd(container),
		cp.CreateCpCmd(container),
		rm.CreateCmd(container),
	)

	middleware.ApplyMiddlewareRecursively(rootCmd, middleware.WithCorrelationID)

	return rootCmd, nil
}
