package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// DefaultService provides the default implementation of the environment service.
type DefaultService struct {
	appName string
}

// NewDefaultService initializes a new instance of the DefaultService.
func NewDefaultService(appName string) *DefaultService {
	return &DefaultService{appName: appName}
}

// OS returns the name of the operating system.
func (s *DefaultService) OS() string {
	return runtime.GOOS
}

// Name returns the name of the application.
func (s *DefaultService) Name() string {
	return s.appName
}

// IsWindows returns true if the OS is Windows.
func (s *DefaultService) IsWindows() bool { return runtime.GOOS == "windows" }

// IsMac returns true if the OS is macOS.
func (s *DefaultService) IsMac() bool { return runtime.GOOS == "darwin" }

// IsLinux returns true if the OS is Linux.
func (s *DefaultService) IsLinux() bool { return runtime.GOOS == "linux" }

// ConfigDir returns the path to the configuration directory.
func (s *DefaultService) ConfigDir() (string, error) {
	if configDir := os.Getenv(EnvConfigDir); strings.TrimSpace(configDir) != "" {
		return configDir, nil
	}

	base, err := configBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, s.appName), nil
}

// DataDir returns the path to the data directory.
func (s *DefaultService) DataDir() (string, error) {
	if dataDir := os.Getenv(EnvDataDir); strings.TrimSpace(dataDir) != "" {
		return dataDir, nil
	}

	base, err := dataBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, s.appName), nil
}

// CacheDir returns the path to the cache directory.
func (s *DefaultService) CacheDir() (string, error) {
	if cacheDir := os.Getenv(EnvCacheDir); strings.TrimSpace(cacheDir) != "" {
		return cacheDir, nil
	}

	base, err := cacheBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, s.appName), nil
}

// LogDir returns the path to the log directory.
func (s *DefaultService) LogDir() (string, error) {
	if logPath := os.Getenv(EnvLogDir); strings.TrimSpace(logPath) != "" {
		return logPath, nil
	}

	if s.IsLinux() {
		base, err := stateBase()
		if err != nil {
			return "", err
		}

		return filepath.Join(base, s.appName, "logs"), nil
	}
	base, err := logsBase()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, s.appName), nil
}

// InstallDir returns the path to the installation directory.
func (s *DefaultService) InstallDir() (string, error) {
	base, err := installBase()
	if err != nil {
		return "", err
	}

	if s.IsWindows() {
		return filepath.Join(base, s.appName), nil
	}

	if s.IsLinux() {
		return base, nil
	}

	return base, nil
}

// TempDir returns the path to the temporary directory.
func (s *DefaultService) TempDir() (string, error) {
	if tempDir := os.Getenv(EnvTempDir); strings.TrimSpace(tempDir) != "" {
		return tempDir, nil
	}

	temp, err := tempBase()
	if err != nil {
		return "", err
	}

	return filepath.Join(temp, s.appName), nil
}

// StateDir returns the path to the state directory.
func (s *DefaultService) StateDir() (string, error) {
	if stateDir := os.Getenv(EnvStateDir); strings.TrimSpace(stateDir) != "" {
		return stateDir, nil
	}

	state, err := stateBase()
	if err != nil {
		return "", err
	}

	return filepath.Join(state, s.appName), nil
}

// EnsureAll ensures all necessary directories exist.
func (s *DefaultService) EnsureAll() error {
	creators := []func() (string, error){
		s.ConfigDir,
		s.DataDir,
		s.CacheDir,
		s.LogDir,
		s.StateDir,
	}

	for _, fn := range creators {
		dir, err := fn()
		if err != nil {
			return err
		}
		if dir == "" {
			return fmt.Errorf("empty directory path resolved")
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	return nil
}

// Shell returns the preferred shell.
func (s *DefaultService) Shell() (string, error) {
	if shell := os.Getenv(EnvShell); strings.TrimSpace(shell) != "" {
		return shell, nil
	}

	return os.Getenv("SHELL"), nil
}

// Editor returns the preferred editor.
func (s *DefaultService) Editor() (string, error) {
	if editor := os.Getenv(EnvEditor); strings.TrimSpace(editor) != "" {
		return editor, nil
	}

	return os.Getenv("EDITOR"), nil
}

// Visual returns the preferred visual editor.
func (s *DefaultService) Visual() (string, error) {
	if visual := os.Getenv(EnvVisual); strings.TrimSpace(visual) != "" {
		return visual, nil
	}

	return os.Getenv("VISUAL"), nil
}
