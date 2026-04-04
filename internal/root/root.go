// Package root provides the entry point for the odc application.
package root

import (
	"fmt"
	"os"
	"strings"

	configcmd "github.com/michaeldcanady/go-onedrive/internal/config/ui/cli"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	drive "github.com/michaeldcanady/go-onedrive/internal/drive/ui/cli"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli/cat"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli/cp"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli/download"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli/edit"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli/ls"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli/mkdir"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli/mv"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli/rm"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli/touch"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli/upload"
	auth "github.com/michaeldcanady/go-onedrive/internal/identity/ui/cli"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/middleware"
	profile "github.com/michaeldcanady/go-onedrive/internal/profile/ui/cli"
	"github.com/michaeldcanady/go-onedrive/internal/shared"
	"github.com/michaeldcanady/go-onedrive/internal/state"
	"github.com/michaeldcanady/go-onedrive/internal/version"
	"github.com/spf13/cobra"
)

const (
	rootShortDescription = "Unix-style OneDrive CLI"
	rootLongDescription  = `odc is a CLI tool designed to interact with OneDrive as a Unix-style file system, providing a terminal-native way to manage files.`
)

// CreateRootCmd constructs and returns the cobra.Command for the root application.
func CreateRootCmd(container di.Container) (*cobra.Command, error) {
	var (
		level       string
		config      string
		profileFlag string
	)

	rootCmd := &cobra.Command{
		Use:           "odc",
		Version:       version.Version,
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
			if exists, _ := container.Profile().Exists(cmd.Context(), shared.DefaultProfileName); !exists {
				_, _ = container.Profile().Create(cmd.Context(), shared.DefaultProfileName)
			}

			if strings.TrimSpace(profileFlag) != "" {
				if err := container.State().Set(state.KeyProfile, profileFlag, state.ScopeSession); err != nil {
					return fmt.Errorf("failed to set session profile: %w", err)
				}
			}

			if strings.TrimSpace(config) != "" {
				if err := container.State().Set(state.KeyConfigOverride, config, state.ScopeSession); err != nil {
					return fmt.Errorf("failed to set config override: %w", err)
				}
			}

			container.Logger().SetAllLevel(level)
			cliLogger.Debug("updated all logger level", logger.String("level", level))
			cliLogger.Debug("updated config path", logger.String("path", config))

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&config, "config", "", "path to config file")
	rootCmd.PersistentFlags().StringVar(&level, "level", "info", "set the logging level (e.g., debug, info, warn, error)")
	rootCmd.PersistentFlags().StringVar(&profileFlag, "profile", "", "name of profile")

	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:
  $ source <(odc completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ odc completion bash > /etc/bash_completion.d/odc
  # macOS:
  $ odc completion bash > /usr/local/etc/bash_completion.d/odc

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ odc completion zsh > "${fpath[1]}/_odc"

  # You will need to start a new shell for this setup to take effect.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				_ = cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				_ = cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				_ = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				_ = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}

	rootCmd.AddCommand(
		auth.CreateAuthCmd(container),
		profile.CreateProfileCmd(container),
		drive.CreateDriveCmd(container),
		configcmd.CreateConfigCmd(container),
		ls.CreateLsCmd(container),
		cat.CreateCatCmd(container),
		mkdir.CreateMkdirCmd(container),
		rm.CreateRmCmd(container),
		touch.CreateTouchCmd(container),
		cp.CreateCpCmd(container),
		mv.CreateMvCmd(container),
		upload.CreateUploadCmd(container),
		download.CreateDownloadCmd(container),
		edit.CreateEditCmd(container),
		completionCmd,
	)

	middleware.ApplyMiddlewareRecursively(rootCmd, middleware.WithCorrelationID)

	return rootCmd, nil
}
