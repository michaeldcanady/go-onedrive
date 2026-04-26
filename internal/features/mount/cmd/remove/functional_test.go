package remove_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/core/env"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	drive "github.com/michaeldcanady/go-onedrive/internal/features/drive/domain"
	editor "github.com/michaeldcanady/go-onedrive/internal/features/editor/domain"
	fsdomain "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount/cmd/remove"
	"github.com/michaeldcanady/go-onedrive/internal/features/profile"
	"github.com/michaeldcanady/go-onedrive/pkg/logger/zap"
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

type testContainer struct {
	mock.Mock
	logger     logger.Service
	config     config.Service
	mounts     mount.Service
	uriFactory *fsdomain.URIFactory
}

func (c *testContainer) Logger() logger.Service           { return c.logger }
func (c *testContainer) Config() config.Service           { return c.config }
func (c *testContainer) Mounts() mount.Service            { return c.mounts }
func (c *testContainer) Identity() identity.Service       { return nil }
func (c *testContainer) Profile() profile.Service         { return nil }
func (c *testContainer) FS() fsdomain.Service             { return nil }
func (c *testContainer) Environment() environment.Service { return nil }
func (c *testContainer) Editor() editor.Service           { return nil }
func (c *testContainer) Drive() drive.Service             { return nil }
func (c *testContainer) URIFactory() *fsdomain.URIFactory { return c.uriFactory }

type mockVFS struct{ mock.Mock; fsdomain.Service }

func (m *mockVFS) Resolve(absPath string) (string, string, error) {
	args := m.Called(absPath)
	return args.String(0), args.String(1), args.Error(2)
}

func TestRemoveCmd_Functional(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	initialConfig := `
mounts:
  - path: "/od:"
    type: "onedrive"
    identity_id: "user1"
`
	err := os.WriteFile(configPath, []byte(initialConfig), 0644)
	require.NoError(t, err)

	envSvc := environment.NewDefaultService("odc-test")
	logSvc := zap.NewZapService(envSvc)
	log, _ := logSvc.CreateLogger("test")

	configSvc := config.NewConfigService(nil, log)
	err = configSvc.SetOverride(context.Background(), configPath)
	require.NoError(t, err)

	mountSvc := mount.NewMountService(&mountConfigAdapter{svc: configSvc})

	mVFS := new(mockVFS)
	mVFS.On("Resolve", "/od").Return("/od", "", nil)

	container := &testContainer{
		logger:     logSvc,
		config:     configSvc,
		mounts:     mountSvc,
		uriFactory: fsdomain.NewURIFactory(mVFS),
	}

	cmd := remove.CreateRemoveCmd(container)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"/od"})

	err = cmd.Execute()
	assert.NoError(t, err)

	mounts, err := mountSvc.ListMounts(context.Background())
	assert.NoError(t, err)
	assert.Len(t, mounts, 0)
}
