// odc is the entry point for the OneDrive CLI, a Unix-style tool for managing
// files across multiple cloud storage providers through a unified virtual namespace.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/core/plugins"
	"github.com/michaeldcanady/go-onedrive/internal/core/resolver"
	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	"github.com/michaeldcanady/go-onedrive/internal/features/drive"
	"github.com/michaeldcanady/go-onedrive/internal/features/editor"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	"github.com/michaeldcanady/go-onedrive/internal/features/profile"
	"github.com/michaeldcanady/go-onedrive/internal/features/storage"
	"github.com/michaeldcanady/go-onedrive/internal/features/vfs"
	"github.com/michaeldcanady/go-onedrive/pkg/format"

	config_get_cmd "github.com/michaeldcanady/go-onedrive/internal/features/config/cmd/config/get"
	config_set_cmd "github.com/michaeldcanady/go-onedrive/internal/features/config/cmd/config/set"

	drive_get_cmd "github.com/michaeldcanady/go-onedrive/internal/features/drive/cmd/drive/get"
	drive_list_cmd "github.com/michaeldcanady/go-onedrive/internal/features/drive/cmd/drive/list"

	identity_list_cmd "github.com/michaeldcanady/go-onedrive/internal/features/identity/cmd/identity/list"
	identity_login_cmd "github.com/michaeldcanady/go-onedrive/internal/features/identity/cmd/identity/login"
	identity_logout_cmd "github.com/michaeldcanady/go-onedrive/internal/features/identity/cmd/identity/logout"

	mount_add_cmd "github.com/michaeldcanady/go-onedrive/internal/features/mount/cmd/mount/add"
	mount_list_cmd "github.com/michaeldcanady/go-onedrive/internal/features/mount/cmd/mount/list"
	mount_remove_cmd "github.com/michaeldcanady/go-onedrive/internal/features/mount/cmd/mount/remove"

	profile_create_cmd "github.com/michaeldcanady/go-onedrive/internal/features/profile/cmd/profile/create"
	profile_current_cmd "github.com/michaeldcanady/go-onedrive/internal/features/profile/cmd/profile/current"
	profile_delete_cmd "github.com/michaeldcanady/go-onedrive/internal/features/profile/cmd/profile/delete"
	profile_list_cmd "github.com/michaeldcanady/go-onedrive/internal/features/profile/cmd/profile/list"
	profile_use_cmd "github.com/michaeldcanady/go-onedrive/internal/features/profile/cmd/profile/use"

	cat_cmd "github.com/michaeldcanady/go-onedrive/internal/features/fs/cmd/cat"
	cp_cmd "github.com/michaeldcanady/go-onedrive/internal/features/fs/cmd/cp"
	download_cmd "github.com/michaeldcanady/go-onedrive/internal/features/fs/cmd/download"
	ls_cmd "github.com/michaeldcanady/go-onedrive/internal/features/fs/cmd/ls"
	mkdir_cmd "github.com/michaeldcanady/go-onedrive/internal/features/fs/cmd/mkdir"
	mv_cmd "github.com/michaeldcanady/go-onedrive/internal/features/fs/cmd/mv"
	rm_cmd "github.com/michaeldcanady/go-onedrive/internal/features/fs/cmd/rm"
	stat_cmd "github.com/michaeldcanady/go-onedrive/internal/features/fs/cmd/stat"
	touch_cmd "github.com/michaeldcanady/go-onedrive/internal/features/fs/cmd/touch"
	upload_cmd "github.com/michaeldcanady/go-onedrive/internal/features/fs/cmd/upload"

	edit_cmd "github.com/michaeldcanady/go-onedrive/internal/features/editor/cmd/edit"
)

var (
	pluginsDir string
	rootCmd    = &cobra.Command{
		Use:     "odc",
		Short:   "OneDrive CLI",
		Version: "0.1.0-dev",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			requestID := uuid.New().String()
			ctx := logger.WithRequestID(cmd.Context(), requestID)
			cmd.SetContext(ctx)
		},
	}
	container di.Container
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic recovered: %v", r)
			if container != nil && container.Logger() != nil {
				container.Logger().Error("application panicked", "error", err)
			} else {
				fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
			}
			os.Exit(1)
		}
	}()

	rootCmd.PersistentFlags().StringVar(&pluginsDir, "plugins-dir", "", "Path to the plugins directory (default: ~/.config/odc/plugins)")

	if err := bootstrap(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := container.Shutdown(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Shutdown error: %v\n", err)
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func bootstrap() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	baseDir := filepath.Join(home, ".config", "odc")
	logDir := filepath.Join(baseDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	l, err := logger.NewZapLogger([]string{filepath.Join(logDir, "app.log"), "stderr"}, "info")
	if err != nil {
		return err
	}

	dbPath := filepath.Join(baseDir, "state.db")
	s, err := storage.NewStorageService(dbPath)
	if err != nil {
		return err
	}

	configRepo := config.NewYAMLRepository(filepath.Join(baseDir, "config.yaml"))
	configService := config.NewConfigService(configRepo, l)

	profileRepo, err := profile.NewBoltRepository(s.DB())
	if err != nil {
		return err
	}
	profileService := profile.NewProfileService(profileRepo, l)

	// Phase 2: Plugins and Identity
	pluginRepo, err := plugins.NewBoltRepository(s.DB())
	if err != nil {
		return err
	}

	if pluginsDir != "" {
		if err := configService.Set(config.KeyCorePluginsDir, pluginsDir); err != nil {
			l.Warn("failed to set plugins directory from flag", "error", err)
		}
	}

	pm := plugins.NewPluginManager(configService, l, pluginRepo)

	identityRepo, err := identity.NewBoltRepository(s.DB())
	if err != nil {
		return err
	}
	ts := identity.NewTokenService(identityRepo, pm, configService, l)
	is := identity.NewIdentityService(identityRepo, pm, ts, l)

	// Phase 3: VFS and Mounts
	mountRepo, err := mount.NewBoltRepository(s.DB())
	if err != nil {
		return err
	}
	ms := mount.NewMountService(mountRepo, l)
	v := vfs.NewOrchestrator(ms, pm, ts, is, l)
	v = vfs.LoggingMiddleware(l)(v)

	// Phase 4: Drive and Editor
	driveRepo, err := drive.NewBoltRepository(s.DB())
	if err != nil {
		return err
	}
	ds := drive.NewDriveService(
		driveRepo,
		pm,
		drive.NewIdentityServiceAdapter(is),
		drive.NewTokenServiceAdapter(ts),
		l,
	)
	es := editor.NewEditorService()

	// Phase 5: Formatting and Resolution
	f := format.NewFactory()
	r := resolver.NewResolverService(v, is, ds, configService)

	container = di.NewContainer(l, s, configService, profileService, pm, ts, is, ms, v, ds, es, f, r)

	registerCommands(container)

	return nil
}

func registerCommands(c di.Container) {
	// Config
	configCmd := &cobra.Command{Use: "config", Short: "Manage configuration"}
	configCmd.AddCommand(config_get_cmd.CreateGetCmd(c))
	configCmd.AddCommand(config_set_cmd.CreateSetCmd(c))
	rootCmd.AddCommand(configCmd)

	// Profile
	profileCmd := &cobra.Command{Use: "profile", Short: "Manage profiles"}
	profileCmd.AddCommand(profile_create_cmd.CreateCreateCmd(c))
	profileCmd.AddCommand(profile_list_cmd.CreateListCmd(c))
	profileCmd.AddCommand(profile_use_cmd.CreateUseCmd(c))
	profileCmd.AddCommand(profile_current_cmd.CreateCurrentCmd(c))
	profileCmd.AddCommand(profile_delete_cmd.CreateDeleteCmd(c))
	rootCmd.AddCommand(profileCmd)

	// Identity
	identityCmd := &cobra.Command{Use: "identity", Short: "Manage identities"}
	identityCmd.AddCommand(identity_login_cmd.CreateLoginCmd(c))
	identityCmd.AddCommand(identity_logout_cmd.CreateLogoutCmd(c))
	identityCmd.AddCommand(identity_list_cmd.CreateListCmd(c))
	rootCmd.AddCommand(identityCmd)

	// Mount
	mountCmd := &cobra.Command{Use: "mount", Short: "Manage mount points"}
	mountCmd.AddCommand(mount_add_cmd.CreateAddCmd(c))
	mountCmd.AddCommand(mount_list_cmd.CreateListCmd(c))
	mountCmd.AddCommand(mount_remove_cmd.CreateRemoveCmd(c))
	rootCmd.AddCommand(mountCmd)

	// Drive
	driveCmd := &cobra.Command{Use: "drive", Short: "Manage storage drives"}
	driveCmd.AddCommand(drive_get_cmd.CreateGetCmd(c))
	driveCmd.AddCommand(drive_list_cmd.CreateListCmd(c))
	rootCmd.AddCommand(driveCmd)

	// FS
	rootCmd.AddCommand(cat_cmd.CreateCatCmd(c))
	rootCmd.AddCommand(cp_cmd.CreateCpCmd(c))
	rootCmd.AddCommand(download_cmd.CreateDownloadCmd(c))
	rootCmd.AddCommand(ls_cmd.CreateLsCmd(c))
	rootCmd.AddCommand(mkdir_cmd.CreateMkdirCmd(c))
	rootCmd.AddCommand(mv_cmd.CreateMvCmd(c))
	rootCmd.AddCommand(rm_cmd.CreateRmCmd(c))
	rootCmd.AddCommand(stat_cmd.CreateStatCmd(c))
	rootCmd.AddCommand(touch_cmd.CreateTouchCmd(c))
	rootCmd.AddCommand(upload_cmd.CreateUploadCmd(c))

	// Editor
	rootCmd.AddCommand(edit_cmd.CreateEditCmd(c))
}
