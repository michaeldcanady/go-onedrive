package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/core/env/infra"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
)

type EnvironmentService struct {
	appName string
}

func New(appName string) *EnvironmentService {
	return &EnvironmentService{appName: appName}
}

func (s *EnvironmentService) OS() string {
	return runtime.GOOS
}

func (s *EnvironmentService) Name() string {
	return s.appName
}

func (s *EnvironmentService) IsWindows() bool { return runtime.GOOS == "windows" }
func (s *EnvironmentService) IsMac() bool     { return runtime.GOOS == "darwin" }
func (s *EnvironmentService) IsLinux() bool   { return runtime.GOOS == "linux" }

func (s *EnvironmentService) ConfigDir() (string, error) {
	if configDir := os.Getenv(infra.EnvConfigDir); strings.TrimSpace(configDir) != "" {
		return configDir, nil
	}

	base, err := infra.ConfigBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, s.appName), nil
}

func (s *EnvironmentService) DataDir() (string, error) {
	if dataDir := os.Getenv(infra.EnvDataDir); strings.TrimSpace(dataDir) != "" {
		return dataDir, nil
	}

	base, err := infra.DataBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, s.appName), nil
}

func (s *EnvironmentService) CacheDir() (string, error) {
	if cacheDir := os.Getenv(infra.EnvCacheDir); strings.TrimSpace(cacheDir) != "" {
		return cacheDir, nil
	}

	base, err := infra.CacheBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, s.appName), nil
}

func (s *EnvironmentService) LogDir() (string, error) {
	if logPath := os.Getenv(infra.EnvLogDir); strings.TrimSpace(logPath) != "" {
		return logPath, nil
	}

	if s.IsLinux() {
		base, err := infra.StateBase()
		if err != nil {
			return "", err
		}

		return filepath.Join(base, s.appName, "logs"), nil
	}
	base, err := infra.LogsBase()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, s.appName), nil
}

func (s *EnvironmentService) InstallDir() (string, error) {
	base, err := infra.InstallBase()
	if err != nil {
		return "", err
	}

	if s.IsWindows() {
		return filepath.Join(base, s.appName), nil
	}

	if s.IsLinux() {
		// ~/.local/bin (binary itself lives here, not a subdirectory)
		return base, nil
	}

	// macOS: ~/Applications (you may later decide on ~/Applications/<app>.app)
	return base, nil
}

func (s *EnvironmentService) TempDir() (string, error) {
	if tempDir := os.Getenv(infra.EnvTempDir); strings.TrimSpace(tempDir) != "" {
		return tempDir, nil
	}

	temp, err := infra.TempBase()
	if err != nil {
		return "", err
	}

	return filepath.Join(temp, s.appName), nil
}

func (s *EnvironmentService) StateDir() (string, error) {
	if stateDir := os.Getenv(infra.EnvStateDir); strings.TrimSpace(stateDir) != "" {
		return stateDir, nil
	}

	state, err := infra.StateBase()
	if err != nil {
		return "", err
	}

	return filepath.Join(state, s.appName), nil
}

func (s *EnvironmentService) EnsureAll() error {
	// Don’t create InstallDir or TempDir
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

func (s *EnvironmentService) OutputDestination() (domainlogger.OutputDestination, error) {
	if outputDest := os.Getenv(infra.EnvLogOutput); strings.TrimSpace(outputDest) != "" {
		if dest := domainlogger.ParseOutputDestination(outputDest); dest != domainlogger.OutputDestinationUnknown {
			return dest, nil
		}
	}

	return domainlogger.DefaultLoggerOutputDestination, nil
}

func (s *EnvironmentService) LogLevel() (string, error) {
	if logLevel := os.Getenv(infra.EnvLogLevel); strings.TrimSpace(logLevel) != "" {
		return logLevel, nil
	}

	return domainlogger.DefaultLoggerLevel, nil
}

func (s *EnvironmentService) Shell() (string, error) {
	// allow for ODC specific shell
	if shell := os.Getenv(infra.EnvShell); strings.TrimSpace(shell) != "" {
		return shell, nil
	}

	return os.Getenv("SHELL"), nil
}

func (s *EnvironmentService) Editor() (string, error) {
	// allow for ODC specific editor
	if editor := os.Getenv(infra.EnvEditor); strings.TrimSpace(editor) != "" {
		return editor, nil
	}

	return os.Getenv("EDITOR"), nil
}

func (s *EnvironmentService) Visual() (string, error) {
	// allow for ODC specific visual editor
	if visual := os.Getenv(infra.EnvVisual); strings.TrimSpace(visual) != "" {
		return visual, nil
	}

	return os.Getenv("VISUAL"), nil
}
