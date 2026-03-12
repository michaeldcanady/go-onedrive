// Package root provides the entry point for the odc2 command-line interface.
package root

import (
	"fmt"
	"os"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/core/shared"
	"github.com/michaeldcanady/go-onedrive/internal/core/state"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/slices/auth"
	"github.com/michaeldcanady/go-onedrive/internal/slices/drive"
	"github.com/michaeldcanady/go-onedrive/internal/slices/drive/cat"
	"github.com/michaeldcanady/go-onedrive/internal/slices/drive/cp"
	"github.com/michaeldcanady/go-onedrive/internal/slices/drive/download"
	"github.com/michaeldcanady/go-onedrive/internal/slices/drive/edit"
	"github.com/michaeldcanady/go-onedrive/internal/slices/drive/ls"
	"github.com/michaeldcanady/go-onedrive/internal/slices/drive/mkdir"
	"github.com/michaeldcanady/go-onedrive/internal/slices/drive/mv"
	"github.com/michaeldcanady/go-onedrive/internal/slices/drive/rm"
	"github.com/michaeldcanady/go-onedrive/internal/slices/drive/touch"
	"github.com/michaeldcanady/go-onedrive/internal/slices/drive/upload"
	"github.com/michaeldcanady/go-onedrive/internal/slices/profile"
	"github.com/spf13/cobra"
)

const (
	rootShortDescription = "A OneDrive CLI client (v2)"
	rootLongDescription  = `A command-line interface for interacting with Microsoft OneDrive, implemented with Vertical Slice Architecture.`
)

// CreateRootCmd constructs and returns the cobra.Command for the root application.
func CreateRootCmd(container di.Container) (*cobra.Command, error) {
	var (
		level       string
		config      string
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
			if exists, _ := container.Profile().Exists(cmd.Context(), shared.DefaultProfileName); !exists {
				_, _ = container.Profile().Create(cmd.Context(), shared.DefaultProfileName)
			}

			if strings.TrimSpace(profileFlag) != "" {
				if err := container.State().Set(state.KeyProfile, profileFlag, state.ScopeSession); err != nil {
					return fmt.Errorf("failed to set session profile: %w", err)
				}
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

	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:
  $ source <(odc2 completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ odc2 completion bash > /etc/bash_completion.d/odc2
  # macOS:
  $ odc2 completion bash > /usr/local/etc/bash_completion.d/odc2

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ odc2 completion zsh > "${fpath[1]}/_odc2"

  # You will need to start a new shell for this setup to take effect.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
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

	return rootCmd, nil
}
