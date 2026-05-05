package profile

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/core/env"
	"github.com/michaeldcanady/go-onedrive/internal/core/shared"
	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

func TestDefaultService_ActiveProfile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "odc-profile-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dbFile := filepath.Join(tmpDir, "profiles.db")
	db, err := bolt.Open(dbFile, 0600, nil)
	assert.NoError(t, err)
	defer db.Close()

	repo, err := NewBoltRepository(db)
	assert.NoError(t, err)

	env := &mockEnv{dir: tmpDir}
	svc := NewDefaultService(repo, repo, env)

	ctx := context.Background()

	// 1. Create profiles
	_, err = svc.Create(ctx, "profile1")
	assert.NoError(t, err)
	_, err = svc.Create(ctx, "profile2")
	assert.NoError(t, err)

	// 2. Set active (global/persistent)
	err = svc.SetActive(ctx, "profile1", shared.ScopeGlobal)
	assert.NoError(t, err)

	active, err := svc.GetActive(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "profile1", active.Name)

	// 3. Set active (session override)
	err = svc.SetActive(ctx, "profile2", shared.ScopeSession)
	assert.NoError(t, err)

	active, err = svc.GetActive(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "profile2", active.Name)

	// 4. Reset session override (implicitly by just checking persistent again if we could)
	// Since we can't easily reset it in this mock without a Reset method,
	// let's just verify that ScopeGlobal still works by setting it again.
	err = svc.SetActive(ctx, "profile1", shared.ScopeGlobal)
	assert.NoError(t, err)
}

type mockEnv struct {
	environment.Service
	dir string
}

func (m *mockEnv) ConfigDir() (string, error) { return m.dir, nil }
func (m *mockEnv) StateDir() (string, error)  { return m.dir, nil }
