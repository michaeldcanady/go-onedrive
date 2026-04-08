package environment

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultService_BasicProperties(t *testing.T) {
	appName := "test-app"
	s := NewDefaultService(appName)

	tests := []struct {
		name     string
		actual   any
		expected any
	}{
		{
			name:     "AppName",
			actual:   s.Name(),
			expected: appName,
		},
		{
			name:     "OS",
			actual:   s.OS(),
			expected: runtime.GOOS,
		},
		{
			name:     "IsWindows",
			actual:   s.IsWindows(),
			expected: runtime.GOOS == "windows",
		},
		{
			name:     "IsMac",
			actual:   s.IsMac(),
			expected: runtime.GOOS == "darwin",
		},
		{
			name:     "IsLinux",
			actual:   s.IsLinux(),
			expected: runtime.GOOS == "linux",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.actual)
		})
	}
}

func TestDefaultService_DirOverrides(t *testing.T) {
	s := NewDefaultService("test-app")

	tests := []struct {
		name     string
		envVar   string
		envVal   string
		call     func() (string, error)
		expected string
	}{
		{
			name:     "ConfigDir override",
			envVar:   EnvConfigDir,
			envVal:   "/custom/config",
			call:     s.ConfigDir,
			expected: "/custom/config",
		},
		{
			name:     "DataDir override",
			envVar:   EnvDataDir,
			envVal:   "/custom/data",
			call:     s.DataDir,
			expected: "/custom/data",
		},
		{
			name:     "CacheDir override",
			envVar:   EnvCacheDir,
			envVal:   "/custom/cache",
			call:     s.CacheDir,
			expected: "/custom/cache",
		},
		{
			name:     "LogDir override",
			envVar:   EnvLogDir,
			envVal:   "/custom/logs",
			call:     s.LogDir,
			expected: "/custom/logs",
		},
		{
			name:     "TempDir override",
			envVar:   EnvTempDir,
			envVal:   "/custom/temp",
			call:     s.TempDir,
			expected: "/custom/temp",
		},
		{
			name:     "StateDir override",
			envVar:   EnvStateDir,
			envVal:   "/custom/state",
			call:     s.StateDir,
			expected: "/custom/state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.envVar, tt.envVal)
			val, err := tt.call()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, val)
		})
	}
}

func TestDefaultService_DefaultPaths(t *testing.T) {
	s := NewDefaultService("test-app")

	tests := []struct {
		name string
		call func() (string, error)
	}{
		{"ConfigDir", s.ConfigDir},
		{"DataDir", s.DataDir},
		{"CacheDir", s.CacheDir},
		{"LogDir", s.LogDir},
		{"TempDir", s.TempDir},
		{"StateDir", s.StateDir},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.call()
			assert.NoError(t, err)
			assert.NotEmpty(t, val)
			assert.Contains(t, val, "test-app")
		})
	}
}

func TestDefaultService_ToolsAndLogs(t *testing.T) {
	s := NewDefaultService("test-app")

	tests := []struct {
		name     string
		envVar   string
		envVal   string
		call     func() string
		expected string
	}{
		{
			name:     "Shell",
			envVar:   EnvShell,
			envVal:   "zsh",
			call:     func() string { v, _ := s.Shell(); return v },
			expected: "zsh",
		},
		{
			name:     "Editor",
			envVar:   EnvEditor,
			envVal:   "vim",
			call:     func() string { v, _ := s.Editor(); return v },
			expected: "vim",
		},
		{
			name:     "Visual",
			envVar:   EnvVisual,
			envVal:   "code",
			call:     func() string { v, _ := s.Visual(); return v },
			expected: "code",
		},
		{
			name:     "LogLevel",
			envVar:   EnvLogLevel,
			envVal:   "debug",
			call:     s.LogLevel,
			expected: "debug",
		},
		{
			name:     "LogOutput",
			envVar:   EnvLogOutput,
			envVal:   "stdout",
			call:     s.LogOutput,
			expected: "stdout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.envVar, tt.envVal)
			assert.Equal(t, tt.expected, tt.call())
		})
	}
}

func TestDefaultService_EnsureAll(t *testing.T) {
	tmpBase := t.TempDir()
	s := NewDefaultService("test-app")

	tests := []struct {
		name string
		setup func()
		check func(t *testing.T)
	}{
		{
			name: "Creates all directories",
			setup: func() {
				t.Setenv(EnvConfigDir, filepath.Join(tmpBase, "config"))
				t.Setenv(EnvDataDir, filepath.Join(tmpBase, "data"))
				t.Setenv(EnvCacheDir, filepath.Join(tmpBase, "cache"))
				t.Setenv(EnvLogDir, filepath.Join(tmpBase, "logs"))
				t.Setenv(EnvStateDir, filepath.Join(tmpBase, "state"))
			},
			check: func(t *testing.T) {
				err := s.EnsureAll()
				assert.NoError(t, err)

				expectedDirs := []string{"config", "data", "cache", "logs", "state"}
				for _, d := range expectedDirs {
					_, err := os.Stat(filepath.Join(tmpBase, d))
					assert.NoError(t, err, "Directory %s should exist", d)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			tt.check(t)
		})
	}
}
