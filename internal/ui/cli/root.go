package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	drive "github.com/michaeldcanady/go-onedrive/internal/drive/ui/cli"
	auth "github.com/michaeldcanady/go-onedrive/internal/identity/ui/cli"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/middleware"

	"github.com/michaeldcanady/go-onedrive/internal/shared"
	"github.com/michaeldcanady/go-onedrive/internal/state"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/config"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/fs/cat"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/fs/cp"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/fs/download"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/fs/edit"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/fs/ls"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/fs/mkdir"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/fs/mv"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/fs/rm"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/fs/touch"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/fs/upload"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/mount"
	"github.com/michaeldcanady/go-onedrive/internal/ui/cli/profile"
	"github.com/michaeldcanady/go-onedrive/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const (
	rootShortDescription = "Unix-style OneDrive CLI"
	rootLongDescription  = `odc is a CLI tool designed to interact with OneDrive as a Unix-style file system, providing a terminal-native way to manage files.`
)

// CreateRootCmd constructs and returns the cobra.Command for the root application.
func CreateRootCmd(container di.Container) (*cobra.Command, error) {
	var (
		levelFlag   string
		configFlag  string
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
			// Ensure a default profile exists
			if exists, _ := container.Profile().Exists(cmd.Context(), shared.DefaultProfileName); !exists {
				_, _ = container.Profile().Create(cmd.Context(), shared.DefaultProfileName)
			}

			if strings.TrimSpace(profileFlag) != "" {
				if err := container.State().Set(state.KeyProfile, profileFlag, state.ScopeSession); err != nil {
					return fmt.Errorf("failed to set session profile: %w", err)
				}
			}

			if strings.TrimSpace(configFlag) != "" {
				if err := container.State().Set(state.KeyConfigOverride, configFlag, state.ScopeSession); err != nil {
					return fmt.Errorf("failed to set config override: %w", err)
				}
			}

			// Load config to get logging settings
			cfg, _ := container.Config().GetConfig(cmd.Context())

			// Determine final log level: CLI flag > Config > Default
			finalLevel := logger.LevelUnknown
			if levelFlag != "" {
				finalLevel = logger.ParseLevel(levelFlag)
				if finalLevel == logger.LevelUnknown {
					return fmt.Errorf("unknown log level: %s", levelFlag)
				}
			} else if cfg.Logging.Level != logger.LevelUnknown {
				finalLevel = cfg.Logging.Level
			}

			// Reconfigure logger with settings from config and flags
			if err := container.Logger().Reconfigure(finalLevel, cfg.Logging.Output, cfg.Logging.Format); err != nil {
				return fmt.Errorf("failed to reconfigure logger: %w", err)
			}

			cliLogger, err := container.Logger().CreateLogger("cli")
			if err != nil {
				return fmt.Errorf("failed to initialize cli logger: %w", err)
			}

			cliLogger.Debug("logger reconfigured",
				logger.String("level", finalLevel.String()),
				logger.String("output", cfg.Logging.Output),
				logger.String("format", cfg.Logging.Format))

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&configFlag, "config", "", "path to config file")
	rootCmd.PersistentFlags().StringVar(&levelFlag, "level", "", "set the logging level (e.g., debug, info, warn, error)")
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

	docsCmd := &cobra.Command{
		Use:    "docs",
		Short:  "Generate documentation",
		Hidden: true,
	}

	manCmd := &cobra.Command{
		Use:   "man [directory]",
		Short: "Generate man pages",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := args[0]
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
			header := &doc.GenManHeader{
				Title:   "ODC",
				Section: "1",
				Source:  "odc " + version.Version,
				Manual:  "odc User Manual",
			}
			return doc.GenManTree(cmd.Root(), header, dir)
		},
	}

	markdownCmd := &cobra.Command{
		Use:   "markdown [directory]",
		Short: "Generate Markdown documentation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := args[0]
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
			return doc.GenMarkdownTree(cmd.Root(), dir)
		},
	}

	docsCmd.AddCommand(manCmd, markdownCmd)

	rootCmd.AddCommand(
		auth.CreateAuthCmd(container),
		profile.CreateProfileCmd(container),
		drive.CreateDriveCmd(container),
		config.CreateConfigCmd(container),
		mount.CreateMountCmd(container),
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
		docsCmd,
	)

	middleware.ApplyMiddlewareRecursively(rootCmd, middleware.WithCorrelationID)

	return rootCmd, nil
}
