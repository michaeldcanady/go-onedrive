package mount_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mountConfigAdapter struct {
	svc config.Service
}

func (a *mountConfigAdapter) GetMounts(ctx context.Context) ([]mount.MountConfig, error) {
	cfg, err := a.svc.GetConfig(ctx)
	if err != nil {
		return nil, err
	}
	var mounts []mount.MountConfig
	for _, m := range cfg.Mounts {
		mounts = append(mounts, mount.MountConfig{
			Path:       m.Path,
			Type:       m.Type,
			IdentityID: m.IdentityID,
			Options:    m.Options,
		})
	}
	return mounts, nil
}

func (a *mountConfigAdapter) SaveMounts(ctx context.Context, mounts []mount.MountConfig) error {
	cfg, err := a.svc.GetConfig(ctx)
	if err != nil {
		return err
	}
	var newMounts []config.MountConfig
	for _, m := range mounts {
		newMounts = append(newMounts, config.MountConfig{
			Path:       m.Path,
			Type:       m.Type,
			IdentityID: m.IdentityID,
			Options:    m.Options,
		})
	}
	cfg.Mounts = newMounts
	return a.svc.SaveConfig(ctx, cfg)
}

type dummyLogger struct{ mock.Mock }

func (l *dummyLogger) Debug(msg string, fields ...logger.Field) {}
func (l *dummyLogger) Info(msg string, fields ...logger.Field)  {}
func (l *dummyLogger) Warn(msg string, fields ...logger.Field)  {}
func (l *dummyLogger) Error(msg string, fields ...logger.Field) {}
func (l *dummyLogger) SetLevel(level logger.Level)             {}
func (l *dummyLogger) With(fields ...logger.Field) logger.Logger { return l }
func (l *dummyLogger) WithContext(ctx context.Context) logger.Logger { return l }

func TestMountFeature_Functional(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create initial empty config
	err := os.WriteFile(configPath, []byte("mounts: []"), 0644)
	require.NoError(t, err)

	ctx := context.Background()
	log := &dummyLogger{}
	configSvc := config.NewConfigService(nil, log)
	err = configSvc.SetOverride(ctx, configPath)
	require.NoError(t, err)

	mountSvc := mount.NewMountService(&mountConfigAdapter{svc: configSvc})

	t.Run("add and then list mount", func(t *testing.T) {
		newCfg := mount.MountConfig{
			Path:       "/od",
			Type:       "onedrive",
			IdentityID: "user1",
			Options:    map[string]string{"drive_id": "123"},
		}

		err := mountSvc.AddMount(ctx, newCfg)
		assert.NoError(t, err)

		mounts, err := mountSvc.ListMounts(ctx)
		assert.NoError(t, err)
		assert.Len(t, mounts, 1)
		assert.Equal(t, "/od", mounts[0].Path)
		assert.Equal(t, "user1", mounts[0].IdentityID)
		assert.Equal(t, "123", mounts[0].Options["drive_id"])
	})

	t.Run("remove mount", func(t *testing.T) {
		err := mountSvc.RemoveMount(ctx, "/od")
		assert.NoError(t, err)

		mounts, err := mountSvc.ListMounts(ctx)
		assert.NoError(t, err)
		assert.Len(t, mounts, 0)
	})
}
